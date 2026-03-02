package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type WorkType string
type ExperienceLevel string
type JobSource string
type JobStatus string

const (
	WorkTypeFullTime   WorkType = "full-time"
	WorkTypePartTime   WorkType = "part-time"
	WorkTypeContract   WorkType = "contract"
	WorkTypeInternship WorkType = "internship"
	WorkTypeRemote     WorkType = "remote"
)

const (
	ExpEntry     ExperienceLevel = "entry"
	ExpMid       ExperienceLevel = "mid"
	ExpSenior    ExperienceLevel = "senior"
	ExpExecutive ExperienceLevel = "executive"
)

const (
	SourceLinkedIn JobSource = "linkedin"
	SourceTwitter  JobSource = "twitter"
	SourceFacebook JobSource = "facebook"
	SourceManual   JobSource = "manual"
)

const (
	StatusPending  JobStatus = "pending"
	StatusApproved JobStatus = "approved"
	StatusRejected JobStatus = "rejected"
)

func (w WorkType) Value() (driver.Value, error)        { return string(w), nil }
func (e ExperienceLevel) Value() (driver.Value, error) { return string(e), nil }
func (s JobSource) Value() (driver.Value, error)       { return string(s), nil }
func (j JobStatus) Value() (driver.Value, error)       { return string(j), nil }

func (w *WorkType) Scan(value interface{}) error {
	if value == nil {
		*w = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*w = WorkType(v)
	case string:
		*w = WorkType(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (e *ExperienceLevel) Scan(value interface{}) error {
	if value == nil {
		*e = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*e = ExperienceLevel(v)
	case string:
		*e = ExperienceLevel(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (s *JobSource) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*s = JobSource(v)
	case string:
		*s = JobSource(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

func (j *JobStatus) Scan(value interface{}) error {
	if value == nil {
		*j = ""
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = JobStatus(v)
	case string:
		*j = JobStatus(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

type Job struct {
	ID              uint            `json:"id" gorm:"primaryKey"`
	Title           string          `json:"title" gorm:"not null;size:255"`
	Company         string          `json:"company" gorm:"not null;size:255"`
	Description     string          `json:"description" gorm:"type:text"`
	Requirements    string          `json:"requirements" gorm:"type:text"`
	Qualifications  string          `json:"qualifications" gorm:"type:text"`
	SalaryMin       *float64        `json:"salary_min"`
	SalaryMax       *float64        `json:"salary_max"`
	SalaryCurrency  string          `json:"salary_currency" gorm:"default:KES;size:10"`
	Location        string          `json:"location" gorm:"size:255"`
	County          string          `json:"county" gorm:"size:100;index"`
	WorkType        WorkType        `json:"work_type" gorm:"type:enum('full-time','part-time','contract','internship','remote');index"`
	ExperienceLevel ExperienceLevel `json:"experience_level" gorm:"type:enum('entry','mid','senior','executive');index"`
	Industry        string          `json:"industry" gorm:"size:100;index"`
	Source          JobSource       `json:"source" gorm:"type:enum('linkedin','twitter','facebook','manual')"`
	SourceURL       string          `json:"source_url" gorm:"size:500;uniqueIndex"`
	Status          JobStatus       `json:"status" gorm:"type:enum('pending','approved','rejected');default:pending;index"`
	RejectionReason string          `json:"rejection_reason" gorm:"type:text"`
	PostedAt        *time.Time      `json:"posted_at"`
	ScrapedAt       *time.Time      `json:"scraped_at"`
	ApprovedAt      *time.Time      `json:"approved_at"`
	ApprovedBy      *uint           `json:"approved_by"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
