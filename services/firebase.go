package services

import (
	"context"
	"log"
	"time"

	"goldprice-api/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	CollectionName     = "gold_prices"
	MetadataCollection = "metadata"
)

var (
	client *firestore.Client
	ctx    = context.Background()
)

// InitFirebase initializes the Firebase Admin SDK
func InitFirebase(credentialsPath string) error {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		return err
	}

	log.Println("✅ Firebase Firestore initialized successfully")
	return nil
}

// CloseFirebase closes the Firestore client
func CloseFirebase() {
	if client != nil {
		client.Close()
	}
}

// GetClient returns the Firestore client
func GetClient() *firestore.Client {
	return client
}

// StoreGoldPrices stores gold price records in Firestore using batch write
// Only adds new records or updates if price has changed
func StoreGoldPrices(prices []models.GoldPrice) (int, error) {
	newCount := 0
	updatedCount := 0
	skippedCount := 0

	// Process in batches of 500 (Firestore limit)
	batchSize := 500
	for i := 0; i < len(prices); i += batchSize {
		end := i + batchSize
		if end > len(prices) {
			end = len(prices)
		}

		batch := client.Batch()
		batchCount := 0

		for _, price := range prices[i:end] {
			// Use sortable date as document ID for proper ordering (2025-12-26)
			docID := price.DateSortable
			docRef := client.Collection(CollectionName).Doc(docID)

			// Check if document exists
			existingDoc, err := docRef.Get(ctx)
			if err != nil || !existingDoc.Exists() {
				// Document doesn't exist, add it
				price.UploadedAt = time.Now()
				batch.Set(docRef, price)
				batchCount++
				newCount++
			} else {
				// Document exists, check if price changed
				var existing models.GoldPrice
				if err := existingDoc.DataTo(&existing); err == nil {
					// Compare prices to see if update is needed
					if existing.Price24K != price.Price24K ||
						existing.Price22K != price.Price22K ||
						existing.Price18K != price.Price18K ||
						existing.PricePerOunce != price.PricePerOunce {
						// Price changed, update it
						price.UploadedAt = time.Now()
						batch.Set(docRef, price)
						batchCount++
						updatedCount++
					} else {
						// Same price, skip
						skippedCount++
					}
				}
			}
		}

		// Commit batch if there are changes
		if batchCount > 0 {
			if _, err := batch.Commit(ctx); err != nil {
				return newCount + updatedCount, err
			}
		}
	}

	log.Printf("📊 Upload summary: %d new, %d updated, %d skipped (no change)",
		newCount, updatedCount, skippedCount)

	return newCount + updatedCount, nil
}

// UpdateMetadata updates the upload metadata
func UpdateMetadata(filename string, recordCount int) error {
	docRef := client.Collection(MetadataCollection).Doc("upload_info")

	metadata := models.Metadata{
		LastUpload:   time.Now(),
		TotalRecords: recordCount,
		Filename:     filename,
	}

	_, err := docRef.Set(ctx, metadata)
	return err
}

// GetAllPrices retrieves all gold prices from Firestore
func GetAllPrices(limit int) ([]models.GoldPrice, error) {
	query := client.Collection(CollectionName).
		OrderBy("date_sortable", firestore.Desc).
		Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var prices []models.GoldPrice
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var price models.GoldPrice
		if err := doc.DataTo(&price); err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}

	return prices, nil
}

// GetLatestPrice retrieves the most recent gold price
func GetLatestPrice() (*models.GoldPrice, error) {
	query := client.Collection(CollectionName).
		OrderBy("date_sortable", firestore.Desc).
		Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var price models.GoldPrice
	if err := doc.DataTo(&price); err != nil {
		return nil, err
	}

	return &price, nil
}

// GetPriceByDate retrieves gold price for a specific date
// Date can be in format "26-Dec-25" or "2025-12-26"
func GetPriceByDate(date string) (*models.GoldPrice, error) {
	// Try to parse and convert to sortable format
	docID := convertToSortableDate(date)
	docRef := client.Collection(CollectionName).Doc(docID)

	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	var price models.GoldPrice
	if err := doc.DataTo(&price); err != nil {
		return nil, err
	}

	return &price, nil
}

// convertToSortableDate converts date from "26-Dec-25" to "2025-12-26"
func convertToSortableDate(dateStr string) string {
	// If already in YYYY-MM-DD format, return as is
	if len(dateStr) == 10 && dateStr[4] == '-' && dateStr[7] == '-' {
		return dateStr
	}

	// Try to parse common formats
	formats := []string{
		"02-Jan-06",
		"02-Jan-2006",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	// If parsing fails, try sanitizing and return
	return sanitizeDateForDocID(dateStr)
}

// GetMetadata retrieves upload metadata
func GetMetadata() (*models.Metadata, error) {
	docRef := client.Collection(MetadataCollection).Doc("upload_info")

	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	var metadata models.Metadata
	if err := doc.DataTo(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

// ClearAllPrices deletes all price data from Firestore
func ClearAllPrices() (int, error) {
	batchSize := 100
	deleted := 0

	for {
		iter := client.Collection(CollectionName).Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return deleted, err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			break
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return deleted, err
		}

		deleted += numDeleted
	}

	return deleted, nil
}

// sanitizeDateForDocID converts date format to Firestore document ID
// Example: "24-Dec-25" -> "24_Dec_25"
func sanitizeDateForDocID(date string) string {
	result := ""
	for _, c := range date {
		if c == '-' {
			result += "_"
		} else {
			result += string(c)
		}
	}
	return result
}
