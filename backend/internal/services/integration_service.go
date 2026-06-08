package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

// IntegrationService handles calls to external modules
type IntegrationService struct {
	client       *http.Client
	inventoryURL string
	financeURL   string
	crmURL       string
	maxRetries   int
}

func NewIntegrationService() *IntegrationService {
	return &IntegrationService{
		client:       &http.Client{Timeout: 10 * time.Second},
		inventoryURL: "http://localhost:8081/api",
		financeURL:   "http://localhost:8082/api",
		crmURL:       "http://localhost:8083/api",
		maxRetries:   3,
	}
}

// CheckInventory checks stock availability with retry logic
func (s *IntegrationService) CheckInventory(productID int, quantity int) (bool, error) {
	url := fmt.Sprintf("%s/inventory/check?product_id=%d&quantity=%d", s.inventoryURL, productID, quantity)
	resp, err := s.retryRequest("GET", url, nil)
	if err != nil {
		log.Printf("Inventory check failed (assuming available): %v", err)
		return true, nil // Graceful degradation
	}
	defer resp.Body.Close()
	var result struct{ Available bool `json:"available"` }
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Available, nil
}

// PushRevenueToFinance sends revenue data to the Finance module
func (s *IntegrationService) PushRevenueToFinance(data map[string]interface{}) error {
	body, _ := json.Marshal(data)
	_, err := s.retryRequest("POST", s.financeURL+"/revenue", body)
	if err != nil {
		log.Printf("Finance push failed: %v", err)
		return err
	}
	return nil
}

// SyncCustomerToCRM syncs customer data to CRM module
func (s *IntegrationService) SyncCustomerToCRM(customerID int) error {
	url := fmt.Sprintf("%s/customers/sync/%d", s.crmURL, customerID)
	_, err := s.retryRequest("POST", url, nil)
	if err != nil {
		log.Printf("CRM sync failed: %v", err)
		return err
	}
	return nil
}

// retryRequest performs an HTTP request with exponential backoff retry
func (s *IntegrationService) retryRequest(method, url string, body []byte) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt < s.maxRetries; attempt++ {
		var req *http.Request
		var err error
		if body != nil {
			req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}
		if err != nil { return nil, err }
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		lastErr = fmt.Errorf("attempt %d failed: %v", attempt+1, err)
		backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
		log.Printf("Retry %d/%d for %s %s (backoff: %v)", attempt+1, s.maxRetries, method, url, backoff)
		time.Sleep(backoff)
	}
	return nil, fmt.Errorf("all %d retries failed: %w", s.maxRetries, lastErr)
}
