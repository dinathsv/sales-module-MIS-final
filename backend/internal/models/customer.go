package models

import "time"

// Customer represents a customer record
type Customer struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	Company    string    `json:"company"`
	TotalSales int       `json:"total_sales,omitempty"`
	TotalSpent float64   `json:"total_spent,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CustomerHistoryFilter holds filtering parameters for customer sales history
type CustomerHistoryFilter struct {
	DateFrom  string
	DateTo    string
	ProductID int
	Status    string
	Page      int
	Limit     int
}
