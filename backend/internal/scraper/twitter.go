package scraper

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"

	"kazi/backend/internal/models"
)

var twitterSearchQueries = []string{
	"hiring Kenya -retweet",
	"job opening Nairobi",
	"vacancy Kenya apply",
	"we are hiring Kenya",
}

func ScrapeTwitter() ([]ScrapedJob, error) {
	var allJobs []ScrapedJob

	browser, cleanup := getBrowser()
	defer cleanup()

	for _, q := range twitterSearchQueries[:2] {
		jobs, err := scrapeTwitterQuery(browser, q)
		if err != nil {
			log.Printf("Twitter query '%s' error: %v", q, err)
			continue
		}
		allJobs = append(allJobs, jobs...)
		time.Sleep(5 * time.Second)
	}

	return allJobs, nil
}

func scrapeTwitterQuery(browser *rod.Browser, query string) ([]ScrapedJob, error) {
	page := browser.MustPage("")
	defer page.MustClose()

	searchURL := fmt.Sprintf(
		"https://x.com/search?q=%s&src=typed_query&f=live",
		strings.ReplaceAll(query, " ", "%20"),
	)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}

	page.WaitLoad()
	time.Sleep(4 * time.Second)

	// Detect login wall — X redirects unauthenticated users to /i/flow/login
	currentURL := page.MustInfo().URL
	if strings.Contains(currentURL, "login") || strings.Contains(currentURL, "i/flow") {
		return nil, fmt.Errorf("X.com requires login — attach an authenticated Edge session via CDP (CDP_URL in .env)")
	}

	var jobs []ScrapedJob

	tweets, err := page.Elements("article[data-testid='tweet']")
	if err != nil || len(tweets) == 0 {
		return jobs, nil
	}

	for i, tweet := range tweets {
		if i >= 8 {
			break
		}

		textEl, err := tweet.Element("[data-testid='tweetText']")
		if err != nil {
			continue
		}
		text := strings.TrimSpace(textEl.MustText())

		if !isJobTweet(text) {
			continue
		}

		tweetURL := ""
		if linkEl, err := tweet.Element("a[href*='/status/']"); err == nil {
			if href, err := linkEl.Attribute("href"); err == nil && href != nil {
				if !strings.HasPrefix(*href, "http") {
					tweetURL = "https://x.com" + *href
				} else {
					tweetURL = *href
				}
			}
		}

		if tweetURL == "" {
			continue
		}

		company := "Unknown"
		if authorEl, err := tweet.Element("[data-testid='User-Name'] span"); err == nil {
			company = strings.TrimSpace(authorEl.MustText())
		}

		title := extractJobTitle(text)
		if title == "" {
			title = "Job Opportunity"
		}

		job := ScrapedJob{
			Title:           title,
			Company:         company,
			Description:     text,
			Source:          models.SourceTwitter,
			SourceURL:       tweetURL,
			Location:        extractLocation(text),
			Industry:        inferIndustry(text),
			WorkType:        inferWorkType(text),
			ExperienceLevel: inferExperienceLevel(text),
		}
		job.County = inferCounty(job.Location)

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func isJobTweet(text string) bool {
	text = strings.ToLower(text)
	jobKeywords := []string{
		"hiring", "job opening", "vacancy", "we are looking for",
		"join our team", "apply now", "open position", "job opportunity",
		"careers", "recruitment", "immediate vacancy",
	}
	for _, kw := range jobKeywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func extractJobTitle(text string) string {
	triggers := []string{"hiring ", "looking for a ", "looking for an ", "seeking a ", "seeking an ", "vacancy: ", "position: "}
	lower := strings.ToLower(text)
	for _, trigger := range triggers {
		idx := strings.Index(lower, trigger)
		if idx >= 0 {
			start := idx + len(trigger)
			rest := text[start:]
			end := strings.IndexAny(rest, "\n.,!?")
			if end < 0 || end > 60 {
				end = 60
			}
			if end > len(rest) {
				end = len(rest)
			}
			title := strings.TrimSpace(rest[:end])
			if len(title) > 5 {
				return title
			}
		}
	}
	return ""
}

func extractLocation(text string) string {
	kenyanCities := []string{
		"Nairobi", "Mombasa", "Kisumu", "Nakuru", "Eldoret", "Thika",
		"Machakos", "Meru", "Nyeri", "Kakamega", "Kisii", "Garissa",
		"Malindi", "Kilifi", "Embu", "Kericho",
	}
	for _, city := range kenyanCities {
		if strings.Contains(text, city) {
			return city + ", Kenya"
		}
	}
	if strings.Contains(strings.ToLower(text), "kenya") {
		return "Kenya"
	}
	return "Kenya"
}
