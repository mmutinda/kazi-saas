package scraper

import (
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kazi/backend/internal/database"
	"kazi/backend/internal/models"
)

type ScrapedJob struct {
	Title           string
	Company         string
	Description     string
	Requirements    string
	Location        string
	County          string
	WorkType        models.WorkType
	ExperienceLevel models.ExperienceLevel
	Industry        string
	Source          models.JobSource
	SourceURL       string
	PostedAt        *time.Time
	SalaryMin       *float64
	SalaryMax       *float64
}

// getBrowser connects to an existing browser via CDP if one is available
// (e.g. Edge launched with --remote-debugging-port=9222), otherwise falls
// back to launching a new headless Chrome instance.
//
// The caller must call cleanup() when done — it is a no-op when we attached
// to an existing browser so we never close a session we don't own.
func getBrowser() (browser *rod.Browser, cleanup func()) {
	cdpURL := os.Getenv("CDP_URL")
	if cdpURL == "" {
		cdpURL = "http://127.0.0.1:9222"
	}

	wsURL, err := launcher.ResolveURL(cdpURL)
	if err == nil {
		log.Printf("Scraper: attaching to existing browser at %s", cdpURL)
		b := rod.New().ControlURL(wsURL).MustConnect()
		// Do NOT close a browser we did not launch.
		return b, func() {}
	}

	log.Printf("Scraper: no CDP browser at %s (%v) — launching headless Chrome", cdpURL, err)
	l := launcher.New().
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu", "true").
		Set("no-first-run", "true").
		MustLaunch()
	b := rod.New().ControlURL(l).MustConnect()
	return b, func() { b.MustClose() }
}

func SaveJobs(jobs []ScrapedJob, source models.JobSource, logID uint) (int, error) {
	saved := 0
	now := time.Now()

	for _, sj := range jobs {
		if sj.SourceURL == "" {
			continue
		}

		job := models.Job{
			Title:           sj.Title,
			Company:         sj.Company,
			Description:     sj.Description,
			Requirements:    sj.Requirements,
			Location:        sj.Location,
			County:          sj.County,
			WorkType:        sj.WorkType,
			ExperienceLevel: sj.ExperienceLevel,
			Industry:        sj.Industry,
			Source:          source,
			SourceURL:       sj.SourceURL,
			Status:          models.StatusPending,
			SalaryCurrency:  "KES",
			SalaryMin:       sj.SalaryMin,
			SalaryMax:       sj.SalaryMax,
			ScrapedAt:       &now,
			PostedAt:        sj.PostedAt,
		}

		result := database.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "source_url"}},
			DoNothing: true,
		}).Create(&job)

		if result.Error != nil {
			log.Printf("Failed to save job %s: %v", sj.SourceURL, result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			saved++
		}
	}

	return saved, nil
}

func startLog(source models.JobSource) uint {
	sl := models.ScrapeLog{
		Source:    source,
		Status:    models.ScrapeRunning,
		StartedAt: time.Now(),
	}
	database.DB.Create(&sl)
	return sl.ID
}

func finishLog(logID uint, found, saved int, errMsg string) {
	now := time.Now()
	status := models.ScrapeCompleted
	if errMsg != "" {
		status = models.ScrapeFailed
	}
	database.DB.Model(&models.ScrapeLog{}).Where("id = ?", logID).Updates(map[string]interface{}{
		"status":        status,
		"jobs_found":    found,
		"jobs_saved":    saved,
		"error_message": errMsg,
		"completed_at":  &now,
	})
}

func RunScraper(source models.JobSource) {
	logID := startLog(source)
	log.Printf("Starting scrape: %s (log #%d)", source, logID)

	var jobs []ScrapedJob
	var err error

	switch source {
	case models.SourceLinkedIn:
		jobs, err = ScrapeLinkedIn()
	case models.SourceTwitter:
		jobs, err = ScrapeTwitter()
	case models.SourceFacebook:
		jobs, err = ScrapeFacebook()
	default:
		finishLog(logID, 0, 0, "unknown source: "+string(source))
		return
	}

	if err != nil {
		log.Printf("Scrape error (%s): %v", source, err)
		finishLog(logID, 0, 0, err.Error())
		return
	}

	// Enrich each job via Ollama (cleans fields, extracts salary, filters non-jobs).
	// Falls back gracefully when Ollama is not reachable.
	enriched := make([]ScrapedJob, 0, len(jobs))
	for _, j := range jobs {
		j, isJob := enrichJob(j)
		if isJob {
			enriched = append(enriched, j)
		}
	}
	log.Printf("Scrape enrich (%s): %d → %d jobs after filtering", source, len(jobs), len(enriched))
	jobs = enriched

	saved, _ := SaveJobs(jobs, source, logID)
	finishLog(logID, len(jobs), saved, "")
	log.Printf("Scrape complete (%s): found=%d saved=%d", source, len(jobs), saved)
}

func StartScheduler() {
	c := cron.New()
	c.AddFunc("0 */6 * * *", func() {
		for _, source := range []models.JobSource{
			models.SourceLinkedIn,
			models.SourceTwitter,
			models.SourceFacebook,
		} {
			go RunScraper(source)
		}
	})
	c.Start()
	log.Println("Scraper scheduler started (every 6h)")

	var count int64
	database.DB.Model(&models.Job{}).Count(&count)
	log.Printf("Existing jobs in DB: %d", count)
}

func GetDB() *gorm.DB {
	return database.DB
}
