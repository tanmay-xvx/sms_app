package transport

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sms-app-backend/common"
)

// HTTPHandler handles HTTP requests for the SMS service
type HTTPHandler struct {
	endpoints Endpoints
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(svc interface{}) *HTTPHandler {
	return &HTTPHandler{
		endpoints: MakeEndpoints(svc),
	}
}

// RegisterRoutes registers all SMS service routes
func (h *HTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	sms := router.Group("/sms")
	{
		sms.POST("/send-otp", h.endpoints.SendOTP)
		sms.POST("/verify-otp", h.endpoints.VerifyOTP)
		sms.POST("/send-sms", h.endpoints.SendSMS)
		sms.GET("/otp-status/:phone", h.endpoints.GetOTPStatus)
	}
	
	callback := router.Group("/callback")
	{
		callback.POST("/request", h.endpoints.RequestCallback)
		callback.GET("/status/:request_id", h.endpoints.GetCallbackStatus)
	}
	
	logs := router.Group("/logs")
	{
		logs.GET("", h.endpoints.GetLogs)
	}
}

// HealthCheck handles health check requests
func (h *HTTPHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "sms-service",
		"version": "1.0.0",
	})
}

// CORSMiddleware handles CORS for the SMS service
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ErrorHandler handles errors and converts them to appropriate HTTP responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// Try to convert to AppError
			if appErr, ok := err.(*common.AppError); ok {
				c.JSON(appErr.StatusCode, appErr)
				return
			}
			
			// Default error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    common.ErrCodeInternal,
				"message": "Internal Server Error",
				"details": err.Error(),
			})
		}
	}
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	// Simple in-memory rate limiter
	// In production, use Redis or similar
	requests := make(map[string][]int64)
	
	return func(c *gin.Context) {
		phone := c.Param("phone")
		if phone == "" {
			// Try to get from request body for POST requests
			if c.Request.Method == "POST" {
				var req struct {
					PhoneNumber string `json:"phone_number"`
				}
				if err := c.ShouldBindJSON(&req); err == nil {
					phone = req.PhoneNumber
				}
			}
		}
		
		if phone != "" {
			now := time.Now().Unix()
			window := now - 60 // 1 minute window
			
			// Clean old requests
			if timestamps, exists := requests[phone]; exists {
				var valid []int64
				for _, ts := range timestamps {
					if ts > window {
						valid = append(valid, ts)
					}
				}
				requests[phone] = valid
				
				// Check rate limit (max 5 requests per minute)
				if len(valid) >= 5 {
					c.JSON(http.StatusTooManyRequests, gin.H{
						"code":    common.ErrCodeRateLimit,
						"message": "Rate limit exceeded",
						"details": "Too many requests. Please try again later.",
					})
					c.Abort()
					return
				}
			}
			
			// Add current request
			requests[phone] = append(requests[phone], now)
		}
		
		c.Next()
	}
} 