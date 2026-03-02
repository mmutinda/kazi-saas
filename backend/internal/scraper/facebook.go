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

var facebookPages = []string{
	"https://www.facebook.com/groups/jobvack",
	"https://web.facebook.com/kenya.new.jobs.vacancies",
	"https://web.facebook.com/profile.php?id=100070331760559",
}

// jsExtractFacebookPosts extracts post text + permalink from a Facebook page/group.
// Works for both authenticated pages and groups in the modern FB UI.
const jsExtractFacebookPosts = `() => {
	const posts = [];

	// Article/post containers — FB uses role="article" or specific feed data attrs
	const articleSelectors = [
		'div[role="feed"] div[role="article"]',
		'div[data-pagelet*="Feed"] div[role="article"]',
		'div[role="article"]',
	];

	let articles = [];
	for (const sel of articleSelectors) {
		articles = Array.from(document.querySelectorAll(sel));
		if (articles.length > 0) break;
	}

	articles.slice(0, 15).forEach((art, idx) => {
		// Text content — try several known containers
		let text = '';
		for (const sel of [
			'div[data-ad-comet-preview="message"]',
			'div[data-ad-preview="message"]',
			'div[dir="auto"]',
			'[data-testid="post_message"]',
		]) {
			const el = art.querySelector(sel);
			if (el && el.innerText && el.innerText.trim().length > 30) {
				text = el.innerText.trim();
				break;
			}
		}

		if (!text) return;

		// Permalink
		let url = '';
		for (const sel of [
			'a[href*="/posts/"]',
			'a[href*="/permalink/"]',
			'a[href*="story_fbid"]',
			'a[role="link"][href*="facebook.com"]',
		]) {
			const el = art.querySelector(sel);
			if (el && el.href && !el.href.includes('/comments/')) {
				url = el.href.split('?')[0];
				break;
			}
		}

		// Fallback URL using post index so we still get a unique key
		if (!url) url = window.location.href + '#post-' + idx;

		posts.push({ text, url });
	});

	return posts;
}`

type fbPost struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func ScrapeFacebook() ([]ScrapedJob, error) {
	var allJobs []ScrapedJob
	seen := map[string]bool{}

	browser, cleanup := getBrowser()
	defer cleanup()

	for _, pageURL := range facebookPages {
		jobs, err := scrapeFacebookPage(browser, pageURL)
		if err != nil {
			log.Printf("Facebook '%s' error: %v", pageURL, err)
			continue
		}
		for _, j := range jobs {
			if !seen[j.SourceURL] {
				seen[j.SourceURL] = true
				allJobs = append(allJobs, j)
			}
		}
		time.Sleep(5 * time.Second)
	}

	return allJobs, nil
}

func scrapeFacebookPage(browser *rod.Browser, pageURL string) ([]ScrapedJob, error) {
	page := browser.MustPage("")
	defer page.MustClose()

	if err := page.Navigate(pageURL); err != nil {
		return nil, fmt.Errorf("navigate: %w", err)
	}
	page.WaitLoad()
	time.Sleep(5 * time.Second)

	// Scroll to load more posts
	page.Mouse.Scroll(0, 800, 3)
	time.Sleep(2 * time.Second)

	info, _ := page.Info()
	log.Printf("Facebook page: %s", info.URL)

	// Detect login wall
	if strings.Contains(info.URL, "login") || strings.Contains(info.URL, "checkpoint") {
		return nil, fmt.Errorf("Facebook requires login — please log in via the Edge CDP browser")
	}

	result, err := page.Eval(jsExtractFacebookPosts)
	if err != nil {
		log.Printf("Facebook JS eval error: %v", err)
		return nil, nil
	}

	raw, _ := json.Marshal(result.Value)
	var posts []fbPost
	if err := json.Unmarshal(raw, &posts); err != nil {
		return nil, fmt.Errorf("unmarshal posts: %w", err)
	}

	log.Printf("Facebook [%s]: found %d posts", pageURL, len(posts))

	pageName := extractFacebookPageName(pageURL)
	var jobs []ScrapedJob

	for _, p := range posts {
		if !isJobPost(p.Text) {
			continue
		}

		title := extractJobTitle(p.Text)
		if title == "" {
			title = "Job Opportunity"
		}

		job := ScrapedJob{
			Title:           title,
			Company:         pageName,
			Description:     p.Text,
			Source:          models.SourceFacebook,
			SourceURL:       p.URL,
			Location:        extractLocation(p.Text),
			Industry:        inferIndustry(p.Text),
			WorkType:        inferWorkType(p.Text),
			ExperienceLevel: inferExperienceLevel(p.Text),
		}
		job.County = inferCounty(job.Location)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func isJobPost(text string) bool {
	text = strings.ToLower(text)
	keywords := []string{
		"hiring", "vacancy", "vacancies", "job opening", "we are looking for",
		"join our team", "apply now", "open position", "job opportunity",
		"recruitment", "immediate vacancy", "career opportunity", "urgently needed",
		"applications invited", "call for applications",
	}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

func extractFacebookPageName(url string) string {
	parts := strings.Split(strings.TrimRight(url, "/"), "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// Strip query string
		name = strings.SplitN(name, "?", 2)[0]
		name = strings.ReplaceAll(name, "-", " ")
		name = strings.ReplaceAll(name, "_", " ")
		words := strings.Fields(name)
		for i, w := range words {
			if len(w) > 0 {
				words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
			}
		}
		return strings.Join(words, " ")
	}
	return "Facebook Jobs Kenya"
}
