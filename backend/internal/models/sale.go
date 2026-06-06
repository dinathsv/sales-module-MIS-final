package models

import "time"

// Sale status constants
const (
	SaleStatusPending   = "pending"
	SaleStatusCompleted = "completed"
	SaleStatusCancelled = "cancelled"
)

// Sale represents a sales transaction
type Sale struct {
	ID              int        `json:"id"`
	CustomerID      *int       `json:"customer_id"`
	CustomerName    string     `json:"customer_name,omitempty"`
	OrderID         string     `json:"order_id"`
	Status          string     `json:"status"`
	Subtotal        float64    `json:"subtotal"`
	DiscountPercent float64    `json:"discount_percent"`
	DiscountAmount  float64    `json:"discount_amount"`
	TaxAmount       float64    `json:"tax_amount"`
	TotalAmount     float64    `json:"total_amount"`
	Notes           string     `json:"notes"`
	CreatedBy       *int       `json:"created_by"`
	Items           []SaleItem `json:"items,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// SaleItem represents a line item in a sale
type SaleItem struct {
	ID        int     `json:"id"`
	SaleID    int     `json:"sale_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	LineTotal float64 `json:"line_total"`
}

// CreateSaleRequest is the request body for creating a sale
type CreateSaleRequest struct {
	CustomerID      int                  `json:"customer_id" binding:"required"`
	DiscountPercent float64              `json:"discount_percent"`
	Notes           string               `json:"notes"`
	Items           []CreateSaleItemRequest `json:"items" binding:"required,min=1"`
}

// CreateSaleItemRequest represents a line item in a create sale request
type CreateSaleItemRequest struct {
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"required"`
}

// UpdateSaleRequest is the request body for updating a sale
type UpdateSaleRequest struct {
	CustomerID      *int                    `json:"customer_id"`
	DiscountPercent *float64                `json:"discount_percent"`
	Notes           *string                 `json:"notes"`
	Items           []CreateSaleItemRequest `json:"items"`
}

// UpdateStatusRequest is the request body for updating sale status
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// SaleFilter holds filtering parameters for listing sales
type SaleFilter struct {
	Status     string
	CustomerID int
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}

// PaginatedResponse wraps paginated results
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
