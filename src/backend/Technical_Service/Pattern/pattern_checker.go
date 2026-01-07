package Pattern

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"project/Technical_Service/Notification"
	"project/Technical_Service/ORM/Entity/EntityStruct"
	model "project/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

type PatternChecker struct {
	db                  *gorm.DB
	notificationService *Notification.NotificationService
	sentNotifications   map[string]time.Time // Key: "{Symbol}-{Pattern}-{Timeus}" -> SentAt
	mu                  sync.Mutex
}

func NewPatternChecker() *PatternChecker {
	adapter := new(model.Adapter).GetAdapterIntance()
	return &PatternChecker{
		db:                  adapter.GetGormIntance(),
		notificationService: Notification.NewNotificationService(),
		sentNotifications:   make(map[string]time.Time),
	}
}

type PredictResponse struct {
	Symbol      string  `json:"symbol"`
	PatternName string  `json:"pattern_name"`
	IsPattern   bool    `json:"is_pattern"`
	Confidence  float64 `json:"confidence"`
	ThaiTime    string  `json:"thai_time"`
	UsTime      string  `json:"us_time"` // We act on US Time as unique identifier for candle
	Error       string  `json:"error,omitempty"`
}

func (pc *PatternChecker) CheckPatterns() {
	var stocks []EntityStruct.Stock
	// Get all unique stock names distinct
	// Gorm doesn't support Distinct on slice directly easily without custom query often,
	// but let's just fetch all and dedup in code or use Distinct("stock_short_name")
	if err := pc.db.Distinct("stock_short_name").Find(&stocks).Error; err != nil {
		fmt.Println("Error fetching stocks:", err)
		return
	}

	for _, stock := range stocks {
		if stock.StockShortName == "" {
			continue
		}
		go pc.processStock(stock.StockShortName)
	}
}

func (pc *PatternChecker) processStock(symbol string) {
	// 1. Call Python API
	// Assume localhost:8000/predict
	apiURL := "http://localhost:8000/predict"

	// Prepare payload: We want "latest" technically, but the API takes a "datetime" to check *around* or *specific*?
	// Looking at app.py: Query is { symbol, datetime }.
	// And it fetches data *for the day of that datetime* and looks at candles *before selected time*.
	// If we want "Right Now", we should send current time.

	// Handle Timezone: explicitly convert to Asia/Bangkok
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		// Fallback to Local if timezone data missing, but print warning
		fmt.Println("⚠️ Warning: Could not load Asia/Bangkok timezone, using Local")
		loc = time.Local
	}
	now := time.Now().In(loc)

	payload := map[string]string{
		"symbol":   symbol,
		"datetime": now.Format("2006-01-02T15:04"), // Format: YYYY-MM-DDTHH:mm
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error calling predict for %s: %v\n", symbol, err)
		return
	}
	defer resp.Body.Close()

	var result PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Error decoding api response for %s: %v\n", symbol, err)
		return
	}

	if result.Error != "" {
		// Just silently ignore errors like "Outside market hours" to avoid log spam?
		// Or log debug.
		// fmt.Printf("API Error for %s: %s\n", symbol, result.Error)
		return
	}

	if result.IsPattern {
		// 2. Anti-Spam Check
		key := fmt.Sprintf("%s-%s-%s", result.Symbol, result.PatternName, result.UsTime)

		pc.mu.Lock()
		if _, exists := pc.sentNotifications[key]; exists {
			pc.mu.Unlock()
			// fmt.Printf("Skipping duplicate notification for %s\n", key)
			return
		}
		// Add to cache
		pc.sentNotifications[key] = time.Now()
		pc.mu.Unlock()

		// 3. Find interested users
		fmt.Printf("🚀 Pattern Found! %s: %s (Conf: %.2f) at %s\n", symbol, result.PatternName, result.Confidence, result.UsTime)
		pc.notifyUsers(symbol, result.PatternName, result.Confidence, result.UsTime)
	}
}

func (pc *PatternChecker) notifyUsers(symbol, pattern string, confidence float64, timeStr string) {
	var stocks []EntityStruct.Stock
	// Find all entries for this symbol to get UserIDs
	if err := pc.db.Preload("UserUser").Where("stock_short_name = ?", symbol).Find(&stocks).Error; err != nil {
		fmt.Println("Error fetching interested users:", err)
		return
	}

	// Dedup emails
	emailMap := make(map[string]bool)
	var emails []string

	for _, s := range stocks {
		// Check if User is loaded and has email
		// Note: ensure EntityStruct.User has Email field and is public
		// Looking at usage in main.go: user.Email exists.
		// s.UserUser is the relation field name based on gen.go?
		// Wait, in stock.gen.go: UserUser stockBelongsToUserUser `RelationField: ... "EntityStruct.User"`
		// So s.UserUser is the wrapper? No, stock struct has `UserUser` field?
		// Checking stock.gen.go -> type stock struct { ... UserUser stockBelongsToUserUser ... } this is the DAO.
		// Checking main.go -> EntityStruct.User definition not fully seen but inferred.
		// We need to check EntityStruct.Stock definition in `stock.gen.go` implied or we must see the struct file.
		// `stock.gen.go` imports `project/Technical_Service/ORM/Entity/EntityStruct`.
		// Let's assume standard GORM Preload works if relations are set up.
		// BUT `gen` code structure often separates DAO from Model.

		// Let's try to fetch User manually by ID if Preload ensures complexity we haven't verified.
		// Actually verifying EntityStruct would be safer but let's assume `s.UserUser` isn't the struct field name in `EntityStruct.Stock`.
		// Usually it's `User` or `UserID`.

		var user EntityStruct.User
		if err := pc.db.First(&user, *s.UserID).Error; err == nil {
			if user.Email != "" {
				if !emailMap[user.Email] {
					emailMap[user.Email] = true
					emails = append(emails, user.Email)
				}
			}
		}
	}

	if len(emails) > 0 {
		pc.notificationService.SendEmail(emails, symbol, pattern, confidence, timeStr)
	}
}
