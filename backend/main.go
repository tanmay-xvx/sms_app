package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "sms-app-backend/docs"
	"sms-app-backend/repository/mongo"
	"sms-app-backend/sms_service"
	"sms-app-backend/sms_service/transport"
)

// @title SMS App Backend API
// @version 1.0
// @description Backend API for SMS application with OTP functionality
// @host localhost:8080
// @BasePath /api
// @schemes http https
// @contact.name API Support
// @contact.email support@smsapp.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	corsOrigins := []string{"http://localhost:3000"}
	
	// Add CORS_ORIGIN if set
	if corsOrigin := os.Getenv("CORS_ORIGIN"); corsOrigin != "" {
		corsOrigins = append(corsOrigins, corsOrigin)
	}
	
	// Allow additional origins in production
	if os.Getenv("ENVIRONMENT") == "production" {
		if additionalOrigins := os.Getenv("ADDITIONAL_CORS_ORIGINS"); additionalOrigins != "" {
			origins := strings.Split(additionalOrigins, ",")
			corsOrigins = append(corsOrigins, origins...)
		}
	}
	
	// Remove duplicates and empty strings
	uniqueOrigins := make([]string, 0)
	seen := make(map[string]bool)
	for _, origin := range corsOrigins {
		if origin != "" && !seen[origin] {
			uniqueOrigins = append(uniqueOrigins, origin)
			seen[origin] = true
		}
	}
	
	config.AllowOrigins = uniqueOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	
	// Apply CORS middleware
	r.Use(cors.New(config))
	
	// Log CORS configuration for debugging
	log.Printf("CORS configured with origins: %v", uniqueOrigins)
	log.Printf("Environment: %s", os.Getenv("ENVIRONMENT"))

	// Initialize MongoDB repository
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	
	repo, err := mongo.NewRepository(mongoURI, "sms_app")
	if err != nil {
		log.Printf("Warning: MongoDB not connected: %v", err)
		log.Println("SMS functionality will be limited")
		repo = nil
	}

	// Initialize SMS service components
	var smsClient transport.SMSClient
	plivoAuthID := os.Getenv("PLIVO_AUTH_ID")
	plivoAuthToken := os.Getenv("PLIVO_AUTH_TOKEN")
	plivoFrom := os.Getenv("PLIVO_FROM_NUMBER")
	
	if plivoAuthID != "" && plivoAuthToken != "" && plivoFrom != "" {
		smsClient = transport.NewPlivoClient(plivoAuthID, plivoAuthToken, plivoFrom)
	} else {
		log.Println("Warning: Plivo credentials not configured, using mock client")
		smsClient = transport.NewMockClient("mock")
	}
	
	var smsService sms_service.SMSService
	var callbackService sms_service.CallbackService
	var logsService sms_service.LogsService
	
	if repo != nil {
		smsService = sms_service.NewSMSService(repo, smsClient)
		callbackService = sms_service.NewCallbackService(repo)
		logsService = sms_service.NewLogsService(repo)
	} else {
		log.Println("Warning: Repository not available, SMS service disabled")
	}
	
	// Create a combined service for the HTTP handler
	combinedService := struct {
		sms_service.SMSService
		sms_service.CallbackService
		sms_service.LogsService
	}{
		smsService,
		callbackService,
		logsService,
	}
	
	smsHandler := transport.NewHTTPHandler(combinedService)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "sms-backend",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Messages
		messages := api.Group("/messages")
		{
			messages.GET("/", getMessages)
			messages.POST("/", createMessage)
			messages.GET("/:id", getMessage)
			messages.PUT("/:id", updateMessage)
			messages.DELETE("/:id", deleteMessage)
		}

		// Users
		users := api.Group("/users")
		{
			users.POST("/register", registerUser)
			users.POST("/login", loginUser)
			users.GET("/profile", authMiddleware(), getUserProfile)
		}

		// AI Service integration
		ai := api.Group("/ai")
		{
			ai.POST("/analyze", analyzeMessage)
			ai.POST("/summarize", summarizeMessages)
		}

		// SMS Service endpoints
		if smsService != nil {
			smsHandler.RegisterRoutes(api)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Message handlers
func getMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"messages": []gin.H{},
		"total":    0,
	})
}

func createMessage(c *gin.Context) {
	var message struct {
		Content string `json:"content" binding:"required"`
		To      string `json:"to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      "msg_123",
		"content": message.Content,
		"to":      message.To,
		"status":  "sent",
	})
}

func getMessage(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"content": "Sample message content",
		"status":  "delivered",
	})
}

func updateMessage(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Message updated successfully",
	})
}

func deleteMessage(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Message deleted successfully",
	})
}

// User handlers
func registerUser(c *gin.Context) {
	var user struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    "user_123",
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func loginUser(c *gin.Context) {
	var login struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual authentication
	token := "jwt_token_here"

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func getUserProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":    "user_123",
		"email": "user@example.com",
		"name":  "John Doe",
	})
}

// AI Service handlers
func analyzeMessage(c *gin.Context) {
	var request struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Call AI service
	c.JSON(http.StatusOK, gin.H{
		"analysis": gin.H{
			"sentiment": "positive",
			"intent":    "greeting",
			"confidence": 0.95,
		},
	})
}

func summarizeMessages(c *gin.Context) {
	var request struct {
		Messages []string `json:"messages" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Call AI service
	c.JSON(http.StatusOK, gin.H{
		"summary": "This is a summary of the provided messages.",
		"key_points": []string{
			"Key point 1",
			"Key point 2",
		},
	})
}

// Middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// TODO: Implement JWT validation
		c.Next()
	}
} 