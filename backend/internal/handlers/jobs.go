package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"kazi/backend/internal/database"
	"kazi/backend/internal/models"
)

type JobsHandler struct{}

func (h *JobsHandler) ListJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Job{}).Where("status = ?", models.StatusApproved)

	if county := c.Query("county"); county != "" {
		query = query.Where("county = ?", county)
	}
	if location := c.Query("location"); location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}
	if industry := c.Query("industry"); industry != "" {
		query = query.Where("industry = ?", industry)
	}
	if workType := c.Query("work_type"); workType != "" {
		query = query.Where("work_type = ?", workType)
	}
	if expLevel := c.Query("experience_level"); expLevel != "" {
		query = query.Where("experience_level = ?", expLevel)
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("MATCH(title, company, description) AGAINST(? IN BOOLEAN MODE)", q+"*")
	}

	var total int64
	query.Count(&total)

	var jobs []models.Job
	if err := query.Order("approved_at DESC, created_at DESC").Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (h *JobsHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job
	if err := database.DB.Where("id = ? AND status = ?", id, models.StatusApproved).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}
	c.JSON(http.StatusOK, job)
}
