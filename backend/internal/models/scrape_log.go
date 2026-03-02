package models

import "time"

type ScrapeStatus string

const (
	ScrapeRunning   ScrapeStatus = "running"
	ScrapeCompleted ScrapeStatus = "completed"
	ScrapeFailed    ScrapeStatus = "failed"
)

func (s *ScrapeStatus) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*s = ScrapeStatus(v)
	case string:
		*s = ScrapeStatus(v)
	}
	return nil
}

type ScrapeLog struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Source       JobSource    `json:"source" gorm:"type:enum('linkedin','twitter','facebook','manual')"`
	Status       ScrapeStatus `json:"status" gorm:"type:enum('running','completed','failed')"`
	JobsFound    int          `json:"jobs_found"`
	JobsSaved    int          `json:"jobs_saved"`
	ErrorMessage string       `json:"error_message" gorm:"type:text"`
	StartedAt    time.Time    `json:"started_at"`
	CompletedAt  *time.Time   `json:"completed_at"`
}
