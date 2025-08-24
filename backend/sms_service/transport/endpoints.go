package transport

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"sms-app-backend/common"
	"sms-app-backend/models"
)

// Endpoints holds all the endpoints for the SMS service
type Endpoints struct {
	SendOTP     gin.HandlerFunc
	VerifyOTP   gin.HandlerFunc
	SendSMS     gin.HandlerFunc
	GetOTPStatus gin.HandlerFunc
	RequestCallback gin.HandlerFunc
	GetCallbackStatus gin.HandlerFunc
	GetLogs     gin.HandlerFunc
}

// MakeEndpoints creates endpoints for the SMS service
func MakeEndpoints(svc interface{}) Endpoints {
	return Endpoints{
		SendOTP:     makeSendOTPEndpoint(svc),
		VerifyOTP:   makeVerifyOTPEndpoint(svc),
		SendSMS:     makeSendSMSEndpoint(svc),
		GetOTPStatus: makeGetOTPStatusEndpoint(svc),
		RequestCallback: makeRequestCallbackEndpoint(svc),
		GetCallbackStatus: makeGetCallbackStatusEndpoint(svc),
		GetLogs:     makeGetLogsEndpoint(svc),
	}
}

// @Summary Send OTP
// @Description Generate and send a 6-digit OTP to the specified phone number
// @Tags SMS
// @Accept json
// @Produce json
// @Param request body models.OTPRequest true "OTP Request"
// @Success 200 {object} models.OTPResponse
// @Failure 400 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /sms/send-otp [post]
func makeSendOTPEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.OTPRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			appErr := common.NewValidationError("Invalid request format: " + err.Error())
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate phone number format
		if !isValidPhoneNumber(req.PhoneNumber) {
			appErr := common.NewValidationError("Invalid phone number format")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Send OTP
		smsSvc, ok := svc.(interface{ SendOTP(ctx context.Context, req models.OTPRequest) (*models.OTPResponse, error) })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		response, err := smsSvc.SendOTP(c.Request.Context(), req)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to send OTP: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// In production, don't return the actual OTP in response
		if response.Success {
			response.OTP = "" // Remove OTP from response for security
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Verify OTP
// @Description Verify the OTP sent to the specified phone number
// @Tags SMS
// @Accept json
// @Produce json
// @Param request body models.VerifyOTPRequest true "OTP Verification Request"
// @Success 200 {object} models.VerifyOTPResponse
// @Failure 400 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /sms/verify-otp [post]
func makeVerifyOTPEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.VerifyOTPRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			appErr := common.NewValidationError("Invalid request format: " + err.Error())
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate phone number format
		if !isValidPhoneNumber(req.PhoneNumber) {
			appErr := common.NewValidationError("Invalid phone number format")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate OTP format (6 digits)
		if !isValidOTP(req.OTP) {
			appErr := common.NewValidationError("Invalid OTP format. Must be 6 digits.")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Verify OTP
		smsSvc, ok := svc.(interface{ VerifyOTP(ctx context.Context, req models.VerifyOTPRequest) (*models.VerifyOTPResponse, error) })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		response, err := smsSvc.VerifyOTP(c.Request.Context(), req)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to verify OTP: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Send SMS
// @Description Send a text message to the specified phone number
// @Tags SMS
// @Accept json
// @Produce json
// @Param request body models.SMSRequest true "SMS Request"
// @Success 200 {object} models.SMSResponse
// @Failure 400 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /sms/send-sms [post]
func makeSendSMSEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SMSRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			appErr := common.NewValidationError("Invalid request format: " + err.Error())
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate phone number format
		if !isValidPhoneNumber(req.PhoneNumber) {
			appErr := common.NewValidationError("Invalid phone number format")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate message length
		if len(req.Message) == 0 || len(req.Message) > 160 {
			appErr := common.NewValidationError("Message must be between 1 and 160 characters")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Send SMS
		smsSvc, ok := svc.(interface{ SendSMS(ctx context.Context, req models.SMSRequest) error })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		err := smsSvc.SendSMS(c.Request.Context(), req)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to send SMS: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		c.JSON(http.StatusOK, models.SMSResponse{
			Success:   true,
			Message:   "SMS sent successfully",
			Timestamp: time.Now(),
		})
	}
}

// @Summary Get OTP Status
// @Description Check the status of OTP for a phone number
// @Tags SMS
// @Accept json
// @Produce json
// @Param phone path string true "Phone Number"
// @Success 200 {object} models.OTPStatus
// @Failure 400 {object} common.AppError
// @Router /sms/otp-status/{phone} [get]
func makeGetOTPStatusEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		phoneNumber := c.Param("phone")
		
		if !isValidPhoneNumber(phoneNumber) {
			appErr := common.NewValidationError("Invalid phone number format")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// This would typically check if an OTP exists and its expiry
		// For security reasons, we don't expose OTP details
		c.JSON(http.StatusOK, models.OTPStatus{
			PhoneNumber: phoneNumber,
			HasActiveOTP: false, // In production, check actual status
			Attempts:    0,
		})
	}
}

// isValidPhoneNumber performs basic phone number validation
func isValidPhoneNumber(phone string) bool {
	// Basic validation: should be at least 10 digits and start with +
	if len(phone) < 10 || phone[0] != '+' {
		return false
	}
	
	// Check if all characters after + are digits
	for i := 1; i < len(phone); i++ {
		if phone[i] < '0' || phone[i] > '9' {
			return false
		}
	}
	
	return true
}

// isValidOTP validates OTP format
func isValidOTP(otp string) bool {
	if len(otp) != 6 {
		return false
	}
	
	// Check if all characters are digits
	for _, char := range otp {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	return true
}

// @Summary Request Callback
// @Description Request a callback call to the specified phone number
// @Tags Callback
// @Accept json
// @Produce json
// @Param request body models.CallbackRequest true "Callback Request"
// @Success 200 {object} models.CallbackResponse
// @Failure 400 {object} common.AppError
// @Failure 500 {object} common.AppError
// @Router /callback/request [post]
func makeRequestCallbackEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CallbackRequest
		
		if err := c.ShouldBindJSON(&req); err != nil {
			appErr := common.NewValidationError("Invalid request format: " + err.Error())
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Validate phone number format
		if !isValidPhoneNumber(req.PhoneNumber) {
			appErr := common.NewValidationError("Invalid phone number format")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Request callback
		callbackSvc, ok := svc.(interface{ RequestCallback(ctx context.Context, req models.CallbackRequest) (*models.CallbackResponse, error) })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		response, err := callbackSvc.RequestCallback(c.Request.Context(), req)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to request callback: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Get Callback Status
// @Description Get the status of a callback request
// @Tags Callback
// @Accept json
// @Produce json
// @Param request_id path string true "Callback Request ID"
// @Success 200 {object} models.Callback
// @Failure 400 {object} common.AppError
// @Failure 404 {object} common.AppError
// @Router /callback/status/{request_id} [get]
func makeGetCallbackStatusEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Param("request_id")
		
		if requestID == "" {
			appErr := common.NewValidationError("Request ID is required")
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		// Get callback status
		callbackSvc, ok := svc.(interface{ GetCallbackStatus(ctx context.Context, requestID string) (*models.Callback, error) })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		callback, err := callbackSvc.GetCallbackStatus(c.Request.Context(), requestID)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to get callback status: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		c.JSON(http.StatusOK, callback)
	}
}

// @Summary Get Activity Logs
// @Description Get all OTP and callback activity logs
// @Tags Logs
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of records (default: 100)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} common.AppError
// @Router /logs [get]
func makeGetLogsEndpoint(svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get limit from query parameter, default to 100
		limitStr := c.DefaultQuery("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 100
		}
		
		// Get logs from service
		logsSvc, ok := svc.(interface{ GetLogs(ctx context.Context, limit int) (map[string]interface{}, error) })
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service not available"})
			return
		}
		
		logs, err := logsSvc.GetLogs(c.Request.Context(), limit)
		if err != nil {
			var appErr *common.AppError
			if e, ok := err.(*common.AppError); ok {
				appErr = e
			} else {
				appErr = common.NewInternalError("Failed to get logs: " + err.Error())
			}
			c.JSON(appErr.StatusCode, appErr)
			return
		}

		c.JSON(http.StatusOK, logs)
	}
} 