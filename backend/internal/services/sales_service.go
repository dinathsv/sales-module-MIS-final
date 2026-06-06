package services

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"sales-module/internal/models"
)

// SalesService handles sales business logic
type SalesService struct {
	db *sql.DB
}

// NewSalesService creates a new SalesService
func NewSalesService(db *sql.DB) *SalesService {
	return &SalesService{db: db}
}

// MaxDiscountPercent is the maximum allowed discount percentage
const MaxDiscountPercent = 50.0

// ListSales returns paginated and filtered sales
func (s *SalesService) ListSales(filter models.SaleFilter) (*models.PaginatedResponse, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 20
	}

	// Build WHERE clause
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("s.status = $%d", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.CustomerID > 0 {
		conditions = append(conditions, fmt.Sprintf("s.customer_id = $%d", argIdx))
		args = append(args, filter.CustomerID)
		argIdx++
	}
	if filter.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("s.created_at >= $%d", argIdx))
		args = append(args, filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("s.created_at <= $%d", argIdx))
		args = append(args, filter.DateTo+"T23:59:59")
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM sales s %s", whereClause)
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("count query failed: %w", err)
	}

	// Fetch page
	offset := (filter.Page - 1) * filter.Limit
	query := fmt.Sprintf(`
		SELECT s.id, s.customer_id, COALESCE(c.name, ''), s.order_id, s.status,
		       s.subtotal, s.discount_percent, s.discount_amount, s.tax_amount,
		       s.total_amount, COALESCE(s.notes, ''), s.created_at, s.updated_at
		FROM sales s
		LEFT JOIN customers c ON s.customer_id = c.id
		%s
		ORDER BY s.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, filter.Limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	sales := []models.Sale{}
	for rows.Next() {
		var sale models.Sale
		err := rows.Scan(
			&sale.ID, &sale.CustomerID, &sale.CustomerName, &sale.OrderID, &sale.Status,
			&sale.Subtotal, &sale.DiscountPercent, &sale.DiscountAmount, &sale.TaxAmount,
			&sale.TotalAmount, &sale.Notes, &sale.CreatedAt, &sale.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		sales = append(sales, sale)
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))

	return &models.PaginatedResponse{
		Data:       sales,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetSale returns a single sale with its items
func (s *SalesService) GetSale(id int) (*models.Sale, error) {
	sale := &models.Sale{}
	err := s.db.QueryRow(`
		SELECT s.id, s.customer_id, COALESCE(c.name, ''), s.order_id, s.status,
		       s.subtotal, s.discount_percent, s.discount_amount, s.tax_amount,
		       s.total_amount, COALESCE(s.notes, ''), s.created_at, s.updated_at
		FROM sales s
		LEFT JOIN customers c ON s.customer_id = c.id
		WHERE s.id = $1
	`, id).Scan(
		&sale.ID, &sale.CustomerID, &sale.CustomerName, &sale.OrderID, &sale.Status,
		&sale.Subtotal, &sale.DiscountPercent, &sale.DiscountAmount, &sale.TaxAmount,
		&sale.TotalAmount, &sale.Notes, &sale.CreatedAt, &sale.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sale not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Fetch sale items
	rows, err := s.db.Query(`
		SELECT id, sale_id, quantity, unit_price, line_total
		FROM sale_items
		WHERE sale_id = $1
		ORDER BY id
	`, id)
	if err != nil {
		return nil, fmt.Errorf("items query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.SaleItem
		err := rows.Scan(
			&item.ID, &item.SaleID, &item.Quantity, &item.UnitPrice, &item.LineTotal,
		)
		if err != nil {
			return nil, fmt.Errorf("item scan failed: %w", err)
		}
		sale.Items = append(sale.Items, item)
	}

	return sale, nil
}

// CreateSale creates a new sale with items
func (s *SalesService) CreateSale(req models.CreateSaleRequest, userID int) (*models.Sale, error) {
	// Validate discount threshold
	if req.DiscountPercent > MaxDiscountPercent {
		return nil, fmt.Errorf("discount cannot exceed %.0f%%", MaxDiscountPercent)
	}
	if req.DiscountPercent < 0 {
		return nil, fmt.Errorf("discount cannot be negative")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback()

	// Calculate subtotal from items
	var subtotal float64
	for i := range req.Items {
		req.Items[i].UnitPrice = math.Round(req.Items[i].UnitPrice*100) / 100
		lineTotal := float64(req.Items[i].Quantity) * req.Items[i].UnitPrice
		subtotal += lineTotal
	}

	// Calculate discount and total
	discountAmount := math.Round(subtotal*req.DiscountPercent) / 100
	totalAmount := subtotal - discountAmount

	// Generate order ID
	orderID := generateOrderID(s.db)

	// Insert sale
	var saleID int
	err = tx.QueryRow(`
		INSERT INTO sales (customer_id, order_id, status, subtotal, discount_percent,
		                   discount_amount, tax_amount, total_amount, notes, created_by)
		VALUES ($1, $2, 'pending', $3, $4, $5, 0.00, $6, $7, $8)
		RETURNING id
	`, req.CustomerID, orderID, subtotal, req.DiscountPercent, discountAmount, totalAmount, req.Notes, userID).Scan(&saleID)
	if err != nil {
		return nil, fmt.Errorf("insert sale failed: %w", err)
	}

	// Insert sale items
	for _, item := range req.Items {
		lineTotal := float64(item.Quantity) * item.UnitPrice
		_, err = tx.Exec(`
			INSERT INTO sale_items (sale_id, quantity, unit_price, line_total)
			VALUES ($1, $2, $3, $4)
		`, saleID, item.Quantity, item.UnitPrice, lineTotal)
		if err != nil {
			return nil, fmt.Errorf("insert item failed: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return s.GetSale(saleID)
}

// UpdateSale updates an existing sale
func (s *SalesService) UpdateSale(id int, req models.UpdateSaleRequest) (*models.Sale, error) {
	// Check sale exists and is pending
	existing, err := s.GetSale(id)
	if err != nil {
		return nil, err
	}
	if existing.Status != models.SaleStatusPending {
		return nil, fmt.Errorf("can only update pending sales")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback()

	// Update sale fields
	if req.CustomerID != nil {
		_, err = tx.Exec("UPDATE sales SET customer_id = $1, updated_at = NOW() WHERE id = $2", *req.CustomerID, id)
		if err != nil {
			return nil, err
		}
	}
	if req.Notes != nil {
		_, err = tx.Exec("UPDATE sales SET notes = $1, updated_at = NOW() WHERE id = $2", *req.Notes, id)
		if err != nil {
			return nil, err
		}
	}

	// If items are provided, replace them
	if req.Items != nil && len(req.Items) > 0 {
		_, err = tx.Exec("DELETE FROM sale_items WHERE sale_id = $1", id)
		if err != nil {
			return nil, err
		}

		var subtotal float64
		for _, item := range req.Items {
			lineTotal := float64(item.Quantity) * item.UnitPrice
			subtotal += lineTotal
			_, err = tx.Exec(`
				INSERT INTO sale_items (sale_id, quantity, unit_price, line_total)
				VALUES ($1, $2, $3, $4)
			`, id, item.Quantity, item.UnitPrice, lineTotal)
			if err != nil {
				return nil, err
			}
		}

		discountPercent := existing.DiscountPercent
		if req.DiscountPercent != nil {
			discountPercent = *req.DiscountPercent
		}
		if discountPercent > MaxDiscountPercent {
			return nil, fmt.Errorf("discount cannot exceed %.0f%%", MaxDiscountPercent)
		}

		discountAmount := math.Round(subtotal*discountPercent) / 100
		totalAmount := subtotal - discountAmount

		_, err = tx.Exec(`
			UPDATE sales SET subtotal = $1, discount_percent = $2, discount_amount = $3,
			       total_amount = $4, updated_at = NOW() WHERE id = $5
		`, subtotal, discountPercent, discountAmount, totalAmount, id)
		if err != nil {
			return nil, err
		}
	} else if req.DiscountPercent != nil {
		if *req.DiscountPercent > MaxDiscountPercent {
			return nil, fmt.Errorf("discount cannot exceed %.0f%%", MaxDiscountPercent)
		}
		discountAmount := math.Round(existing.Subtotal*(*req.DiscountPercent)) / 100
		totalAmount := existing.Subtotal - discountAmount
		_, err = tx.Exec(`
			UPDATE sales SET discount_percent = $1, discount_amount = $2,
			       total_amount = $3, updated_at = NOW() WHERE id = $4
		`, *req.DiscountPercent, discountAmount, totalAmount, id)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}

	return s.GetSale(id)
}

// UpdateStatus updates the sale status with validation rules
func (s *SalesService) UpdateStatus(id int, newStatus string) (*models.Sale, error) {
	sale, err := s.GetSale(id)
	if err != nil {
		return nil, err
	}

	// Validate status transitions
	switch sale.Status {
	case models.SaleStatusPending:
		if newStatus != models.SaleStatusCompleted && newStatus != models.SaleStatusCancelled {
			return nil, fmt.Errorf("pending sale can only be completed or cancelled")
		}
	case models.SaleStatusCompleted:
		return nil, fmt.Errorf("completed sale cannot change status")
	case models.SaleStatusCancelled:
		return nil, fmt.Errorf("cancelled sale cannot change status")
	}

	_, err = s.db.Exec("UPDATE sales SET status = $1, updated_at = NOW() WHERE id = $2", newStatus, id)
	if err != nil {
		return nil, fmt.Errorf("update status failed: %w", err)
	}

	return s.GetSale(id)
}

// DeleteSale cancels/deletes a sale (only pending sales)
func (s *SalesService) DeleteSale(id int) error {
	sale, err := s.GetSale(id)
	if err != nil {
		return err
	}
	if sale.Status != models.SaleStatusPending {
		return fmt.Errorf("can only delete pending sales")
	}

	_, err = s.db.Exec("DELETE FROM sales WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}

// generateOrderID generates a unique order ID
func generateOrderID(db *sql.DB) string {
	year := time.Now().Format("2006")
	var count int
	db.QueryRow("SELECT COUNT(*) FROM sales WHERE order_id LIKE $1", fmt.Sprintf("ORD-%s-%%", year)).Scan(&count)
	return fmt.Sprintf("ORD-%s-%04d", year, count+1)
}
