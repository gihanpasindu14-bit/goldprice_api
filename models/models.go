package models

import "time"

// GoldPrice represents a gold price record
type GoldPrice struct {
	Date          string    `firestore:"date" json:"date"`                   // Display format: "26-Dec-25"
	DateSortable  string    `firestore:"date_sortable" json:"date_sortable"` // Sortable format: "2025-12-26"
	PricePerOunce float64   `firestore:"price_per_ounce" json:"price_per_ounce"`
	Price24K      int       `firestore:"price_24k" json:"price_24k"`
	Price22K      int       `firestore:"price_22k" json:"price_22k"`
	Price18K      int       `firestore:"price_18k" json:"price_18k"`
	UploadedAt    time.Time `firestore:"uploaded_at,omitempty" json:"-"`
}

// Metadata represents upload metadata
type Metadata struct {
	LastUpload   time.Time `firestore:"last_upload" json:"last_upload"`
	TotalRecords int       `firestore:"total_records" json:"total_records"`
	Filename     string    `firestore:"filename" json:"filename"`
}

// PriceResponse represents a simplified price response
type PriceResponse struct {
	Date          string  `json:"date"`
	Carat         string  `json:"carat"`
	PriceLKR      int     `json:"price_lkr"`
	PricePerOunce float64 `json:"price_per_ounce"`
}

// LatestPricesResponse represents the latest prices response
type LatestPricesResponse struct {
	Success       bool                  `json:"success"`
	Date          string                `json:"date"`
	Prices        map[string]CaratPrice `json:"prices"`
	PricePerOunce float64               `json:"price_per_ounce"`
}

// CaratPrice represents price info for a specific carat
type CaratPrice struct {
	PriceLKR int    `json:"price_lkr"`
	Carat    string `json:"carat"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
}

// UploadResponse represents CSV upload response
type UploadResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	RecordsProcessed int    `json:"records_processed"`
	Filename         string `json:"filename"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Endpoints map[string]string `json:"endpoints"`
}
