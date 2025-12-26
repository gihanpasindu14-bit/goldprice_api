package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"goldprice-api/models"
	"goldprice-api/services"
	"goldprice-api/utils"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns API status and available endpoints
func HealthCheck(c *gin.Context) {
	response := models.HealthCheckResponse{
		Status:  "online",
		Service: "Gold Price API",
		Version: "1.0.0",
		Endpoints: map[string]string{
			"upload":      "/api/upload (POST)",
			"get_all":     "/api/prices (GET)",
			"get_latest":  "/api/prices/latest (GET)",
			"get_by_date": "/api/prices/{date} (GET)",
			"metadata":    "/api/metadata (GET)",
		},
	}
	c.JSON(http.StatusOK, response)
}

// UploadCSV handles CSV file upload and storage
func UploadCSV(c *gin.Context) {
	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No file uploaded",
		})
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasSuffix(header.Filename, ".csv") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "File must be a CSV",
		})
		return
	}

	// Read file content
	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read file",
		})
		return
	}

	// Validate and parse CSV
	parsedData, err := utils.ValidateAndParseCSV(string(content))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Store data in Firestore
	count, err := services.StoreGoldPrices(parsedData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Failed to store data: %v", err),
		})
		return
	}

	// Update metadata
	if err := services.UpdateMetadata(header.Filename, count); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Warning: Failed to update metadata: %v\n", err)
	}

	message := fmt.Sprintf("Successfully processed %d records (new or updated)", count)
	if count == 0 {
		message = "No new or updated records found. All data is already up to date."
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		Success:          true,
		Message:          message,
		RecordsProcessed: count,
		Filename:         header.Filename,
	})
}

// GetAllPrices retrieves all gold prices with optional filters
func GetAllPrices(c *gin.Context) {
	// Parse query parameters
	limitStr := c.DefaultQuery("limit", "100")
	carat := c.Query("carat")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	// Get prices from Firestore
	prices, err := services.GetAllPrices(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Error fetching data: %v", err),
		})
		return
	}

	// Transform data based on carat filter
	var responseData []models.PriceResponse

	if carat != "" {
		// Filter by specific carat
		caratUpper := strings.ToUpper(carat)
		for _, price := range prices {
			var priceValue int

			switch caratUpper {
			case "24K":
				priceValue = price.Price24K
			case "22K":
				priceValue = price.Price22K
			case "18K":
				priceValue = price.Price18K
			default:
				continue
			}

			responseData = append(responseData, models.PriceResponse{
				Date:          price.Date,
				Carat:         caratUpper,
				PriceLKR:      priceValue,
				PricePerOunce: price.PricePerOunce,
			})
		}
	} else {
		// Return all carats
		for _, price := range prices {
			responseData = append(responseData, models.PriceResponse{
				Date:          price.Date,
				Carat:         "24K",
				PriceLKR:      price.Price24K,
				PricePerOunce: price.PricePerOunce,
			})
			responseData = append(responseData, models.PriceResponse{
				Date:          price.Date,
				Carat:         "22K",
				PriceLKR:      price.Price22K,
				PricePerOunce: price.PricePerOunce,
			})
			responseData = append(responseData, models.PriceResponse{
				Date:          price.Date,
				Carat:         "18K",
				PriceLKR:      price.Price18K,
				PricePerOunce: price.PricePerOunce,
			})
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Count:   len(responseData),
		Data:    responseData,
	})
}

// GetLatestPrices retrieves the most recent gold prices
func GetLatestPrices(c *gin.Context) {
	price, err := services.GetLatestPrice()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Error fetching latest prices: %v", err),
		})
		return
	}

	if price == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "No data available",
		})
		return
	}

	response := models.LatestPricesResponse{
		Success: true,
		Date:    price.Date,
		Prices: map[string]models.CaratPrice{
			"24K": {
				PriceLKR: price.Price24K,
				Carat:    "24K",
			},
			"22K": {
				PriceLKR: price.Price22K,
				Carat:    "22K",
			},
			"18K": {
				PriceLKR: price.Price18K,
				Carat:    "18K",
			},
		},
		PricePerOunce: price.PricePerOunce,
	}

	c.JSON(http.StatusOK, response)
}

// GetPriceByDate retrieves gold prices for a specific date
func GetPriceByDate(c *gin.Context) {
	date := c.Param("date")

	price, err := services.GetPriceByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Error fetching data: %v", err),
		})
		return
	}

	if price == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   fmt.Sprintf("No data found for date: %s", date),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"date":    price.Date,
		"prices": map[string]int{
			"24K": price.Price24K,
			"22K": price.Price22K,
			"18K": price.Price18K,
		},
		"price_per_ounce": price.PricePerOunce,
	})
}

// GetMetadata retrieves upload metadata
func GetMetadata(c *gin.Context) {
	metadata, err := services.GetMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Error fetching metadata: %v", err),
		})
		return
	}

	if metadata == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"metadata": gin.H{
				"message": "No uploads yet",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"metadata": gin.H{
			"last_upload":   metadata.LastUpload,
			"total_records": metadata.TotalRecords,
			"filename":      metadata.Filename,
		},
	})
}

// ClearAllData deletes all price data from Firestore
func ClearAllData(c *gin.Context) {
	deleted, err := services.ClearAllPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Error clearing data: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Deleted %d records", deleted),
	})
}
