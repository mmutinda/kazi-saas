package scraper

import (
	"bytes"
	"encoding/json"
	"html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"kazi/backend/internal/models"
)

const enrichModel = "qwen2.5:7b-instruct"

var reHTMLTag = regexp.MustCompile(`<[^>]+>`)
var reWhitespace = regexp.MustCompile(`\s{2,}`)

// stripHTML removes HTML tags, unescapes entities, and collapses whitespace.
// Mirrors what cheerio's $.text() does: extract readable text from markup.
func stripHTML(s string) string {
	s = reHTMLTag.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = reWhitespace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

const enrichPrompt = `You are a job listing parser for the Kenyan job market.
Extract structured fields from the job post below and return ONLY a valid JSON object with exactly these fields:

{
  "title": "Clean job title",
  "company": "Company or organization name",
  "location": "City or area in Kenya",
  "county": "One of Kenya's 47 counties, or empty string if unknown",
  "work_type": "full-time|part-time|contract|internship|remote",
  "experience_level": "entry|mid|senior|executive",
  "industry": "Technology|Finance|Marketing|Healthcare|Education|Sales|Logistics|HR|Legal|Agriculture|Construction|Hospitality|Media|Other",
  "requirements": "Key requirements or qualifications extracted from the post (empty string if none)",
  "salary_min": 0,
  "salary_max": 0,
  "is_job_post": true
}

Rules:
- salary_min / salary_max: numeric KES values, use 0 if not mentioned
- is_job_post: set to false if this is clearly NOT a job listing (e.g. news, ads, general posts)
- county: map common cities — Nairobi→Nairobi, Mombasa→Mombasa, Kisumu→Kisumu, Nakuru→Nakuru, Eldoret→Uasin Gishu, Thika/Kiambu→Kiambu, Machakos→Machakos
- Return ONLY the JSON object — no markdown fences, no explanation

Job post:
`

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
}

type enrichedFields struct {
	Title           string  `json:"title"`
	Company         string  `json:"company"`
	Location        string  `json:"location"`
	County          string  `json:"county"`
	WorkType        string  `json:"work_type"`
	ExperienceLevel string  `json:"experience_level"`
	Industry        string  `json:"industry"`
	Requirements    string  `json:"requirements"`
	SalaryMin       float64 `json:"salary_min"`
	SalaryMax       float64 `json:"salary_max"`
	IsJobPost       bool    `json:"is_job_post"`
}

// enrichJob calls Ollama (qwen2.5:7b-instruct) to clean and enrich a scraped job's fields.
// HTML is stripped from the description before sending. Falls back gracefully when
// Ollama is not reachable — the job is kept as-is and isJob is true.
func enrichJob(job ScrapedJob) (ScrapedJob, bool) {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	apiURL := ollamaHost + "/api/chat"

	// Strip HTML tags + unescape entities, then truncate.
	text := stripHTML(job.Description)
	if len(text) > 3000 {
		text = text[:3000]
	}
	if text == "" {
		text = job.Title + " at " + job.Company + " — " + job.Location
	}

	reqBody, _ := json.Marshal(ollamaRequest{
		Model:    enrichModel,
		Messages: []ollamaMessage{{Role: "user", Content: enrichPrompt + text}},
		Stream:   false,
	})

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("enrich: build request: %v", err)
		return job, true
	}
	httpReq.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Printf("enrich: Ollama not reachable (%s): %v", ollamaHost, err)
		return job, true
	}
	defer resp.Body.Close()

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil || ollamaResp.Message.Content == "" {
		log.Printf("enrich: decode response: %v", err)
		return job, true
	}

	// Strip accidental markdown fences the model may add.
	raw := strings.TrimSpace(ollamaResp.Message.Content)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var fields enrichedFields
	if err := json.Unmarshal([]byte(raw), &fields); err != nil {
		log.Printf("enrich: parse JSON: %v — raw: %s", err, raw)
		return job, true
	}

	// Model flagged this as not a real job listing.
	if !fields.IsJobPost {
		return job, false
	}

	// Apply enriched values; keep originals when the model returns empty strings.
	if fields.Title != "" {
		job.Title = fields.Title
	}
	if fields.Company != "" {
		job.Company = fields.Company
	}
	if fields.Location != "" {
		job.Location = fields.Location
	}
	if fields.County != "" {
		job.County = fields.County
	}
	if fields.WorkType != "" {
		job.WorkType = models.WorkType(fields.WorkType)
	}
	if fields.ExperienceLevel != "" {
		job.ExperienceLevel = models.ExperienceLevel(fields.ExperienceLevel)
	}
	if fields.Industry != "" {
		job.Industry = fields.Industry
	}
	if fields.Requirements != "" {
		job.Requirements = fields.Requirements
	}
	if fields.SalaryMin > 0 {
		job.SalaryMin = &fields.SalaryMin
	}
	if fields.SalaryMax > 0 {
		job.SalaryMax = &fields.SalaryMax
	}

	return job, true
}
