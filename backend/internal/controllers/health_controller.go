package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController { return &HealthController{} }

func (hc *HealthController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "sales-module", "version": "1.0.0"})
}
