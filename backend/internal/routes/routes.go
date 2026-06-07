package routes

import (
	"database/sql"

	"sales-module/internal/controllers"
	"sales-module/internal/middleware"
	"sales-module/internal/services"

	"github.com/gin-gonic/gin"
)

// Setup registers all routes on the Gin engine
func Setup(r *gin.Engine, db *sql.DB) {
	// Services
	salesService := services.NewSalesService(db)
	reportService := services.NewReportService(db)

	// Controllers
	authCtrl := controllers.NewAuthController(db)
	healthCtrl := controllers.NewHealthController()
	salesCtrl := controllers.NewSalesController(salesService, nil) // passing nil for invoiceService for now
	reportCtrl := controllers.NewReportController(reportService)

	api := r.Group("/api")
	{
		// Public routes
		api.GET("/health", healthCtrl.Health)
		api.POST("/auth/login", authCtrl.Login)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Sales
			protected.GET("/sales", salesCtrl.List)
			protected.POST("/sales", salesCtrl.Create)
			protected.GET("/sales/:id", salesCtrl.Get)
			protected.PUT("/sales/:id", salesCtrl.Update)
			protected.PATCH("/sales/:id/status", salesCtrl.UpdateStatus)
			protected.DELETE("/sales/:id", salesCtrl.Delete)

			// Reports
			protected.GET("/reports/dashboard", reportCtrl.Dashboard)
			protected.GET("/reports/summary", reportCtrl.Summary)
			protected.GET("/reports/revenue", reportCtrl.Revenue)
			protected.POST("/reports/export", reportCtrl.Export)
		}
	}
}
