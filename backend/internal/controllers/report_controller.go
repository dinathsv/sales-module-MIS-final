package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sales-module/internal/services"
)

type ReportController struct {
	service *services.ReportService
}

func NewReportController(s *services.ReportService) *ReportController {
	return &ReportController{service: s}
}

func (rc *ReportController) Dashboard(c *gin.Context) {
	stats, err := rc.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (rc *ReportController) Summary(c *gin.Context) {
	period := c.DefaultQuery("period", "monthly")
	summary, err := rc.service.GetSalesSummary(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// TopProducts removed
func (rc *ReportController) Revenue(c *gin.Context) {
	report, err := rc.service.GetRevenueReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}

func (rc *ReportController) Export(c *gin.Context) {
	// Stub: In production, this would push to Analytics/MIS module
	summary, _ := rc.service.GetSalesSummary("monthly")
	revenue, _ := rc.service.GetRevenueReport()

	c.JSON(http.StatusOK, gin.H{
		"message":      "Data exported successfully",
		"summary":      summary,
		"revenue":      revenue,
	})
}
