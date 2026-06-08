package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sales-module/internal/models"
	"sales-module/internal/services"
)

type SalesController struct {
	service        *services.SalesService
}

func NewSalesController(s *services.SalesService, inv interface{}) *SalesController {
	return &SalesController{service: s}
}

func (sc *SalesController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	custID, _ := strconv.Atoi(c.Query("customer_id"))

	filter := models.SaleFilter{
		Status:     c.Query("status"),
		CustomerID: custID,
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
		Page:       page,
		Limit:      limit,
	}

	result, err := sc.service.ListSales(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (sc *SalesController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}
	sale, err := sc.service.GetSale(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sale)
}

func (sc *SalesController) Create(c *gin.Context) {
	var req models.CreateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("userID")
	sale, err := sc.service.CreateSale(req, userID.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sale)
}

func (sc *SalesController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}
	var req models.UpdateSaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sale, err := sc.service.UpdateSale(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sale)
}

func (sc *SalesController) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sale, err := sc.service.UpdateStatus(id, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Removed auto-generate invoice logic
	// if req.Status == models.SaleStatusCompleted {
	// 	go sc.invoiceService.GenerateInvoice(id)
	// }

	c.JSON(http.StatusOK, sale)
}

func (sc *SalesController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}
	if err := sc.service.DeleteSale(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sale deleted successfully"})
}
