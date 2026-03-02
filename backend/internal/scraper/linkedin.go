package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-rod/rod"

	"kazi/backend/internal/models"
)

var linkedInKeywords = []string{
	"software engineer Kenya",
	"developer Kenya",
	"marketing Kenya",
	"finance Kenya",
	"project manager Kenya",
	"sales Kenya",
	"data analyst Kenya",
	"accounting Kenya",
}

// jsExtractLinkedInJobs runs inside the browser to extract job cards.
// Uses multiple fallback selectors so it survives LinkedIn UI redesigns.
const jsExtractLinkedInJobs = `() => {
	const jobs = [];

	// --- find job card container ---
	const cardSelectors = [
		'li.jobs-search-results__list-item',
		'.jobs-search-results-grid__list li',
		'.scaffold-layout__list-container li',
		'.jobs-search__results-list li',
		'ul.jobs-search__results-list > li',
		'div.job-card-container',
	];
	let cards = [];
	for (const sel of cardSelectors) {
		cards = Array.from(document.querySelectorAll(sel));
		if (cards.length > 0) break;
	}

	cards.slice(0, 15).forEach(card => {
		// --- title ---
		let title = '';
		for (const sel of [
			'.job-card-list__title--link',
			'.job-card-list__title',
			'.job-card-container__link',
			'a[class*="title"]',
			'strong',
			'h3', 'h4',
		]) {
			const el = card.querySelector(sel);
			if (el && el.textContent.trim()) { title = el.textContent.trim(); break; }
		}

		// --- company ---
		let company = '';
		for (const sel of [
			'.job-card-container__company-name',
			'.job-card-container__primary-description',
			'.artdeco-entity-lockup__subtitle span',
			'[class*="company-name"]',
			'[class*="employer"]',
		]) {
			const el = card.querySelector(sel);
			if (el && el.textContent.trim()) { company = el.textContent.trim(); break; }
		}

		// --- location ---
		let location = 'Kenya';
		for (const sel of [
			'.job-card-container__metadata-item--location',
			'.job-card-container__metadata-item',
			'.artdeco-entity-lockup__caption li',
			'[class*="location"]',
			'li[class*="metadata"]',
		]) {
			const el = card.querySelector(sel);
			if (el && el.textContent.trim()) { location = el.textContent.trim(); break; }
		}

		// --- url ---
		let url = '';
		for (const sel of [
			'a.job-card-list__title--link',
			'a.job-card-container__link',
			'a[href*="/jobs/view/"]',
			'a[href*="linkedin.com/jobs"]',
			'a',
		]) {
			const el = card.querySelector(sel);
			if (el && el.href && el.href.includes('/jobs/')) {
				url = el.href.split('?')[0];
				break;
			}
		}

		if (title && url) jobs.push({ title, company, location, url });
	});

	return jobs;
}`

type linkedInCard struct {
	Title    string `json:"title"`
	Company  string `json:"company"`
	Location string `json:"location"`
	URL      string `json:"url"`
}

func ScrapeLinkedIn() ([]ScrapedJob, error) {
	var allJobs []ScrapedJob
	seen := map[string]bool{}

	browser, cleanup := getBrowser()
	defer cleanup()

	for _, keyword := range linkedInKeywords[:4] {
		jobs, err := scrapeLinkedInKeyword(browser, keyword)
		if err != nil {
			log.Printf("LinkedIn keyword '%s' error: %v", keyword, err)
			continue
		}
		for _, j := range jobs {
			if !seen[j.SourceURL] {
				seen[j.SourceURL] = true
				allJobs = append(allJobs, j)
			}
		}
		time.Sleep(3 * time.Second)
	}

	return allJobs, nil
}

// jsCountSelectors counts how many elements each candidate selector matches.
// Used to diagnose why card extraction returns 0.
const jsCountSelectors = `() => {
	const sels = [
		'li.jobs-search-results__list-item',
		'.jobs-search-results-grid__list li',
		'.scaffold-layout__list-container li',
		'.jobs-search__results-list li',
		'ul.jobs-search__results-list > li',
		'div.job-card-container',
	];
	const counts = {};
	for (const s of sels) counts[s] = document.querySelectorAll(s).length;
	counts['_total_li'] = document.querySelectorAll('li').length;
	counts['_total_a_jobs'] = document.querySelectorAll('a[href*="/jobs/"]').length;
	counts['_page_title'] = document.title;
	return counts;
}`

func scrapeLinkedInKeyword(browser *rod.Browser, keyword string) ([]ScrapedJob, error) {
	page := browser.MustPage("")
	defer page.MustClose()

	searchURL := fmt.Sprintf(
		"https://www.linkedin.com/jobs/search/?keywords=%s&location=Kenya&f_TPR=r86400&refresh=true",
		strings.ReplaceAll(keyword, " ", "%20"),
	)
	log.Printf("LinkedIn [%s]: navigating to %s", keyword, searchURL)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}
	page.WaitLoad()
	time.Sleep(3 * time.Second)

	info, _ := page.Info()
	log.Printf("LinkedIn [%s]: landed url=%s", keyword, info.URL)

	// Detect redirect to auth
	if strings.Contains(info.URL, "/authwall") || strings.Contains(info.URL, "/login") || strings.Contains(info.URL, "/checkpoint") {
		return nil, fmt.Errorf("LinkedIn redirected to auth wall — log in via the Edge CDP browser")
	}

	// Run JS extraction
	result, err := page.Eval(jsExtractLinkedInJobs)
	if err != nil {
		log.Printf("LinkedIn [%s]: JS eval error: %v", keyword, err)
		return nil, nil
	}

	raw, err := json.Marshal(result.Value)
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	log.Printf("LinkedIn [%s]: raw JS result: %s", keyword, raw)

	var cards []linkedInCard
	if err := json.Unmarshal(raw, &cards); err != nil {
		return nil, fmt.Errorf("unmarshal cards: %w", err)
	}
	log.Printf("LinkedIn [%s]: %d cards parsed from JS", keyword, len(cards))

	// If nothing was found, run selector diagnostics to show which DOM elements exist.
	if len(cards) == 0 {
		if diag, diagErr := page.Eval(jsCountSelectors); diagErr == nil {
			diagRaw, _ := json.Marshal(diag.Value)
			log.Printf("LinkedIn [%s]: selector diagnostics: %s", keyword, diagRaw)
		}
	}

	var jobs []ScrapedJob
	for _, c := range cards {
		if c.Title == "" {
			log.Printf("LinkedIn [%s]: skipping card — empty title (url=%s)", keyword, c.URL)
			continue
		}
		if c.URL == "" {
			log.Printf("LinkedIn [%s]: skipping card — empty url (title=%q)", keyword, c.Title)
			continue
		}
		job := ScrapedJob{
			Title:           c.Title,
			Company:         c.Company,
			Location:        c.Location,
			Source:          models.SourceLinkedIn,
			SourceURL:       c.URL,
			Industry:        inferIndustry(c.Title + " " + keyword),
			WorkType:        inferWorkType(c.Title),
			ExperienceLevel: inferExperienceLevel(c.Title),
		}
		job.County = inferCounty(job.Location)
		log.Printf("LinkedIn [%s]: queued job title=%q company=%q url=%s", keyword, job.Title, job.Company, job.SourceURL)
		jobs = append(jobs, job)
	}
	log.Printf("LinkedIn [%s]: %d jobs queued for enrichment", keyword, len(jobs))

	return jobs, nil
}

func inferCounty(location string) string {
	loc := strings.ToLower(location)
	countyMap := map[string]string{
		"nairobi":      "Nairobi",
		"mombasa":      "Mombasa",
		"kisumu":       "Kisumu",
		"nakuru":       "Nakuru",
		"eldoret":      "Uasin Gishu",
		"thika":        "Kiambu",
		"kiambu":       "Kiambu",
		"machakos":     "Machakos",
		"meru":         "Meru",
		"nyeri":        "Nyeri",
		"kakamega":     "Kakamega",
		"kisii":        "Kisii",
		"garissa":      "Garissa",
		"malindi":      "Kilifi",
		"kilifi":       "Kilifi",
		"embu":         "Embu",
		"kericho":      "Kericho",
		"bomet":        "Bomet",
		"homabay":      "Homa Bay",
		"homa bay":     "Homa Bay",
		"migori":       "Migori",
		"siaya":        "Siaya",
		"bungoma":      "Bungoma",
		"busia":        "Busia",
		"vihiga":       "Vihiga",
		"trans nzoia":  "Trans Nzoia",
		"west pokot":   "West Pokot",
		"samburu":      "Samburu",
		"laikipia":     "Laikipia",
		"nyandarua":    "Nyandarua",
		"murang'a":     "Murang'a",
		"muranga":      "Murang'a",
		"kirinyaga":    "Kirinyaga",
		"tharaka":      "Tharaka-Nithi",
		"isiolo":       "Isiolo",
		"marsabit":     "Marsabit",
		"mandera":      "Mandera",
		"wajir":        "Wajir",
		"turkana":      "Turkana",
		"baringo":      "Baringo",
		"nandi":        "Nandi",
		"uasin gishu":  "Uasin Gishu",
		"elgeyo":       "Elgeyo-Marakwet",
		"narok":        "Narok",
		"kajiado":      "Kajiado",
		"makueni":      "Makueni",
		"kitui":        "Kitui",
		"kwale":        "Kwale",
		"taita":        "Taita-Taveta",
		"lamu":         "Lamu",
		"tana river":   "Tana River",
	}
	for key, county := range countyMap {
		if strings.Contains(loc, key) {
			return county
		}
	}
	return ""
}

func inferIndustry(text string) string {
	text = strings.ToLower(text)
	industries := map[string][]string{
		"Technology":   {"software", "developer", "engineer", "tech", "it ", "data", "ai ", "machine learning", "devops", "cloud", "cyber", "frontend", "backend", "fullstack"},
		"Finance":      {"finance", "accounting", "accountant", "auditor", "financial", "banking", "bank", "insurance", "investment", "treasury", "tax"},
		"Marketing":    {"marketing", "digital marketing", "seo", "social media", "brand", "advertising", "content", "communications", "pr "},
		"Healthcare":   {"health", "medical", "doctor", "nurse", "clinical", "pharmacy", "hospital", "lab", "dental"},
		"Education":    {"teacher", "education", "school", "tutor", "academic", "lecturer", "training", "curriculum"},
		"Sales":        {"sales", "business development", "account manager", "revenue", "commercial"},
		"Logistics":    {"logistics", "supply chain", "warehouse", "transport", "delivery", "procurement", "fleet"},
		"HR":           {"human resource", "hr ", "recruitment", "talent", "people", "payroll"},
		"Legal":        {"legal", "lawyer", "advocate", "compliance", "counsel", "paralegal"},
		"Agriculture":  {"agriculture", "agri", "farming", "food", "agritech", "horticulture"},
		"Construction": {"construction", "civil", "architect", "building", "quantity survey"},
		"Hospitality":  {"hospitality", "hotel", "tourism", "chef", "restaurant", "catering", "front office"},
		"Media":        {"media", "journalist", "reporter", "communications", "public relations", "videograp", "photography"},
	}
	for industry, keywords := range industries {
		for _, kw := range keywords {
			if strings.Contains(text, kw) {
				return industry
			}
		}
	}
	return "Other"
}

func inferWorkType(text string) models.WorkType {
	text = strings.ToLower(text)
	if strings.Contains(text, "remote") || strings.Contains(text, "wfh") || strings.Contains(text, "work from home") {
		return models.WorkTypeRemote
	}
	if strings.Contains(text, "part-time") || strings.Contains(text, "part time") {
		return models.WorkTypePartTime
	}
	if strings.Contains(text, "contract") || strings.Contains(text, "freelance") || strings.Contains(text, "consultancy") {
		return models.WorkTypeContract
	}
	if strings.Contains(text, "intern") || strings.Contains(text, "attachment") || strings.Contains(text, "trainee") {
		return models.WorkTypeInternship
	}
	return models.WorkTypeFullTime
}

func inferExperienceLevel(title string) models.ExperienceLevel {
	title = strings.ToLower(title)
	if strings.Contains(title, "senior") || strings.Contains(title, "lead") || strings.Contains(title, "principal") || strings.Contains(title, "sr.") || strings.Contains(title, "head of") {
		return models.ExpSenior
	}
	if strings.Contains(title, "junior") || strings.Contains(title, "entry") || strings.Contains(title, "graduate") || strings.Contains(title, "intern") || strings.Contains(title, "trainee") || strings.Contains(title, "attachment") {
		return models.ExpEntry
	}
	if strings.Contains(title, "director") || strings.Contains(title, "vp ") || strings.Contains(title, "chief") || strings.Contains(title, "c-suite") || strings.Contains(title, "president") {
		return models.ExpExecutive
	}
	return models.ExpMid
}
