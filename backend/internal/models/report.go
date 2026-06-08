package models

// SalesSummary holds aggregated sales data for a time period
type SalesSummary struct {
	Period           string  `json:"period"`
	TotalSales       int     `json:"total_sales"`
	TotalRevenue     float64 `json:"total_revenue"`
	CompletedSales   int     `json:"completed_sales"`
	PendingSales     int     `json:"pending_sales"`
	CancelledSales   int     `json:"cancelled_sales"`
	AverageOrderValue float64 `json:"average_order_value"`
}

// RevenueReport holds revenue totals
type RevenueReport struct {
	TotalRevenue     float64          `json:"total_revenue"`
	TotalTransactions int             `json:"total_transactions"`
	GrowthPercent    float64          `json:"growth_percent"`
	RevenueByPeriod  []PeriodRevenue  `json:"revenue_by_period"`
}

// PeriodRevenue holds revenue for a specific period
type PeriodRevenue struct {
	Period  string  `json:"period"`
	Revenue float64 `json:"revenue"`
	Count   int     `json:"count"`
}

// TopProduct holds top-selling product data
type TopProduct struct {
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductSKU   string  `json:"product_sku"`
	Category     string  `json:"category"`
	TotalQtySold int     `json:"total_qty_sold"`
	TotalRevenue float64 `json:"total_revenue"`
}

// DashboardStats holds overview statistics for the dashboard
type DashboardStats struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalSales        int     `json:"total_sales"`
	PendingOrders     int     `json:"pending_orders"`
	AverageOrderValue float64 `json:"average_order_value"`
	RevenueGrowth     float64 `json:"revenue_growth"`
	SalesGrowth       float64 `json:"sales_growth"`
}
