package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kazi/backend/internal/database"
	"kazi/backend/internal/middleware"
	"kazi/backend/internal/models"
	"kazi/backend/internal/scraper"
)

type AdminHandler struct {
	JWTSecret string
}

// POST /api/admin/login
func (h *AdminHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var admin models.Admin
	if err := database.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := &middleware.Claims{
		AdminID: admin.ID,
		Email:   admin.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"admin": gin.H{
			"id":    admin.ID,
			"email": admin.Email,
			"name":  admin.Name,
		},
	})
}

// GET /api/admin/jobs
func (h *AdminHandler) ListJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Job{})

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if source := c.Query("source"); source != "" {
		query = query.Where("source = ?", source)
	}
	if q := c.Query("q"); q != "" {
		query = query.Where("title LIKE ? OR company LIKE ?", "%"+q+"%", "%"+q+"%")
	}

	var total int64
	query.Count(&total)

	var jobs []models.Job
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
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

// GET /api/admin/stats
func (h *AdminHandler) Stats(c *gin.Context) {
	var pending, approved, rejected, total int64
	database.DB.Model(&models.Job{}).Where("status = ?", models.StatusPending).Count(&pending)
	database.DB.Model(&models.Job{}).Where("status = ?", models.StatusApproved).Count(&approved)
	database.DB.Model(&models.Job{}).Where("status = ?", models.StatusRejected).Count(&rejected)
	database.DB.Model(&models.Job{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"pending":  pending,
		"approved": approved,
		"rejected": rejected,
		"total":    total,
	})
}

// PUT /api/admin/jobs/:id/approve
func (h *AdminHandler) ApproveJob(c *gin.Context) {
	id := c.Param("id")
	adminID := c.GetUint("admin_id")
	now := time.Now()

	result := database.DB.Model(&models.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.StatusApproved,
		"approved_at": &now,
		"approved_by": adminID,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve job"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job approved"})
}

// PUT /api/admin/jobs/:id/reject
func (h *AdminHandler) RejectJob(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	result := database.DB.Model(&models.Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":           models.StatusRejected,
		"rejection_reason": req.Reason,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject job"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job rejected"})
}

// DELETE /api/admin/jobs/:id
func (h *AdminHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Where("id = ?", id).Delete(&models.Job{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job deleted"})
}

// POST /api/admin/scrape
func (h *AdminHandler) TriggerScrape(c *gin.Context) {
	var req struct {
		Source string `json:"source" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validSources := map[string]bool{
		"linkedin": true, "twitter": true, "facebook": true,
	}
	if !validSources[req.Source] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid source. Use: linkedin, twitter, facebook"})
		return
	}

	go scraper.RunScraper(models.JobSource(req.Source))

	c.JSON(http.StatusAccepted, gin.H{"message": "Scrape triggered for " + req.Source})
}

// GET /api/admin/scrape/logs
func (h *AdminHandler) ScrapeLogs(c *gin.Context) {
	var logs []models.ScrapeLog
	if err := database.DB.Order("started_at DESC").Limit(50).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GET /api/admin/jobs/:id
func (h *AdminHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	var job models.Job
	if err := database.DB.First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job"})
		return
	}
	c.JSON(http.StatusOK, job)
}
