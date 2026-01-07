package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"project/Technical_Service/ORM/Entity/EntityStruct"
	"project/Technical_Service/Scheduler"
	model "project/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Initialize Database
	adapter := new(model.Adapter).GetAdapterIntance()
	db := adapter.GetGormIntance()

	// Initialize Router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Start Scheduler
	scheduler := Scheduler.NewSchedulerService()
	scheduler.Start()

	// Proxy /predict to Python Service (port 8000)
	r.POST("/predict", func(c *gin.Context) {
		// 1. Read Body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		// Restore body for further binding if needed (not needed here but good practice)
		// c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 2. Forward to Python
		pythonURL := "http://localhost:8000/predict"
		resp, err := http.Post(pythonURL, "application/json", bytes.NewBuffer(bodyBytes))
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Prediction service unavailable"})
			return
		}
		defer resp.Body.Close()

		// 3. Return Response
		c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
	})

	// Routes
	r.POST("/register", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Email    string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		passwordStr := string(hashedPassword)
		usernameStr := input.Username
		emailStr := input.Email
		// Check if user exists
		var existingUser EntityStruct.User
		if result := db.Where("username = ?", input.Username).First(&existingUser); result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		// Create user
		newUser := EntityStruct.User{
			Username: usernameStr,
			Password: passwordStr,
			Email:    emailStr,
		}

		// Generate random UserID (Workaround for missing auto-increment)
		newUser.UserID = int32(rand.Int31())

		if err := db.Create(&newUser).Error; err != nil {
			fmt.Println("❌ Error creating user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	r.POST("/login", func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user EntityStruct.User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if user.Password == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// In a real app, generate JWT here. For now, just return success.
		var usernameVal string
		if user.Username != "" {
			usernameVal = user.Username
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Login successful",
			"userId":   user.UserID,
			"username": usernameVal,
		})
	})

	// Add Stock Endpoint
	r.POST("/add-stock", func(c *gin.Context) {
		var input struct {
			UserID    int32  `json:"userId" binding:"required"`
			StockName string `json:"stockName" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newStock := EntityStruct.Stock{
			UserID:         &input.UserID,
			StockShortName: input.StockName,
			StockID:        int32(rand.Int31()), // Generate random ID
		}

		if err := db.Create(&newStock).Error; err != nil {
			fmt.Println("❌ Error adding stock:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add stock"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully", "stock": newStock})
	})

	// Get User Stocks Endpoint
	r.GET("/user-stocks", func(c *gin.Context) {
		userIdStr := c.Query("userId")
		if userIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}

		var stocks []EntityStruct.Stock
		if err := db.Where("user_id = ?", userIdStr).Find(&stocks).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"stocks": stocks})
	})

	// Run server
	fmt.Println("Server running on port 8080")
	r.Run(":8080")
}
