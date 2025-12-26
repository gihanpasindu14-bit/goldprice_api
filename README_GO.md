# Gold Price API

Go backend for managing gold price data from Sri Lanka Central Bank with Firebase Firestore integration.

## Features

- ✅ Upload CSV files via REST API
- ✅ CSV validation and data parsing
- ✅ Firebase Firestore storage
- ✅ Multiple endpoints for data retrieval
- ✅ CORS enabled for frontend access
- ✅ Automatic price conversion (24K, 22K, 18K)
- 🔜 Phase 2: Auto web scraping from Central Bank

## Setup Instructions

### 1. Install Go

Download and install Go from [https://golang.org/dl/](https://golang.org/dl/) (Go 1.21 or higher)

### 2. Install Dependencies

```bash
cd goldprice_api
go mod download
```

### 3. Configure Firebase

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project: **goldprice-api**
3. Go to **Project Settings** → **Service Accounts**
4. Click **Generate New Private Key**
5. Save the JSON file as `firebase-credentials.json` in this folder
6. Enable Firestore Database in Firebase Console

### 4. Run the Server

```bash
go run main.go
```

Or build and run:

```bash
go build -o goldprice-api
./goldprice-api  # Linux/Mac
goldprice-api.exe  # Windows
```

Server will start at: `http://localhost:8000`

## Project Structure

```
goldprice_api/
├── main.go                      # Entry point, router setup
├── go.mod                       # Go module dependencies
├── firebase-credentials.json    # Firebase service account key
├── handlers/
│   └── handlers.go             # HTTP request handlers
├── services/
│   └── firebase.go             # Firestore operations
├── models/
│   └── models.go               # Data structures
└── utils/
    └── csv_parser.go           # CSV parsing utilities
```

## API Endpoints

### 📋 Health Check
```
GET /
```

### 📤 Upload CSV
```
POST /api/upload
Content-Type: multipart/form-data
Body: file (CSV file)
```

**Example with curl:**
```bash
curl -X POST http://localhost:8000/api/upload \
  -F "file=@../goldprice_lk/data.csv"
```

**Example with Postman:**
1. Select POST method
2. URL: `http://localhost:8000/api/upload`
3. Body → form-data
4. Key: `file` (type: File)
5. Value: Select your CSV file
6. Click Send

### 📊 Get All Prices
```
GET /api/prices?limit=100&carat=24k
```

Query Parameters:
- `limit`: Number of records (default: 100)
- `carat`: Filter by carat (24k, 22k, 18k)

### 🔥 Get Latest Prices
```
GET /api/prices/latest
```

### 📅 Get Price by Date
```
GET /api/prices/24-Dec-25
```

### 📈 Get Metadata
```
GET /api/metadata
```

### 🗑️ Clear All Data (Use with caution!)
```
DELETE /api/prices/clear
```

## CSV Format

Expected format (Central Bank style):
```csv
Header Row 1
Header Row 2
,24-Dec-25,44871.00
,23-Dec-25,44650.00
```

## Response Format

### Upload Response
```json
{
  "success": true,
  "message": "Successfully uploaded 230 records",
  "records_processed": 230,
  "filename": "data.csv"
}
```

### Get Prices Response
```json
{
  "success": true,
  "count": 3,
  "data": [
    {
      "date": "24-Dec-25",
      "carat": "24K",
      "price_lkr": 44871
    },
    {
      "date": "24-Dec-25",
      "carat": "22K",
      "price_lkr": 41132
    }
  ]
}
```

### Latest Prices Response
```json
{
  "success": true,
  "date": "24-Dec-25",
  "prices": {
    "24K": {
      "price_lkr": 44871,
      "carat": "24K"
    },
    "22K": {
      "price_lkr": 41132,
      "carat": "22K"
    },
    "18K": {
      "price_lkr": 33653,
      "carat": "18K"
    }
  },
  "price_per_ounce": 1395567.89
}
```

## Firestore Collections

- `gold_prices`: Stores all price records (document ID = date)
- `metadata`: Stores upload metadata

## Security Notes

⚠️ **Current Setup (Development):**
- CORS: Allow all origins (`*`)
- No API key required
- Public access

🔒 **Phase 2 (Production):**
- Add API key authentication
- Restrict CORS to specific domains
- Add rate limiting
- Environment-based configuration

## Testing

### Test Upload
```bash
curl -X POST http://localhost:8000/api/upload \
  -F "file=@../goldprice_lk/data.csv"
```

### Test Get Latest
```bash
curl http://localhost:8000/api/prices/latest
```

## Troubleshooting

**Error: "cannot find package"**
- Solution: `go mod download`

**Error: "Could not open firebase-credentials.json"**
- Solution: Download service account key from Firebase Console

**Error: "Permission denied" on Firestore**
- Solution: Check Firestore security rules, set to test mode initially

## Next Steps (Phase 2)

- [ ] Web scraping from Central Bank website
- [ ] Automated CSV download scheduler
- [ ] API key authentication
- [ ] Rate limiting
- [ ] Deploy to cloud (Google Cloud Run)
- [ ] Add caching layer
- [ ] Email notifications on new data

## Tech Stack

- **Framework**: Gin (Go web framework)
- **Database**: Firebase Firestore
- **Language**: Go 1.21+
- **Server**: Built-in Go HTTP server

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `firebase.google.com/go/v4` - Firebase Admin SDK
- `cloud.google.com/go/firestore` - Firestore client

## License

MIT
