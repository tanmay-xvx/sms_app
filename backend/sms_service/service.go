package sms_service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"sms-app-backend/common"
	"sms-app-backend/models"
	"sms-app-backend/repository"
	"sms-app-backend/sms_service/transport"
)

// SMSServiceImpl implements the SMSService interface
type SMSServiceImpl struct {
	repo        repository.Repository
	smsClient   transport.SMSClient
}

// CallbackServiceImpl implements the CallbackService interface
type CallbackServiceImpl struct {
	repo repository.Repository
}

// LogsServiceImpl implements the LogsService interface
type LogsServiceImpl struct {
	repo repository.Repository
}

// NewSMSService creates a new SMS service instance
func NewSMSService(repo repository.Repository, smsClient transport.SMSClient) *SMSServiceImpl {
	service := &SMSServiceImpl{
		repo:      repo,
		smsClient: smsClient,
	}

	// Start cleanup goroutine
	go service.startCleanupRoutine()

	return service
}

// SendSMS sends a regular SMS message
func (s *SMSServiceImpl) SendSMS(ctx context.Context, req models.SMSRequest) error {
	log.Printf("Sending SMS to %s: %s", req.PhoneNumber, req.Message)
	
	// Create SMS record
	sms := &models.SMS{
		From:     s.smsClient.GetProvider(),
		To:       req.PhoneNumber,
		Message:  req.Message,
		Status:   models.StatusPending,
		Provider: s.smsClient.GetProvider(),
	}

	// Store SMS record
	err := s.repo.SMS().Create(ctx, sms)
	if err != nil {
		log.Printf("Failed to store SMS record: %v", err)
		return common.NewInternalError("Failed to store SMS record")
	}

	// Send SMS via provider
	err = s.smsClient.SendSMS(ctx, req.PhoneNumber, req.Message)
	if err != nil {
		log.Printf("Failed to send SMS to %s: %v", req.PhoneNumber, err)
		
		// Update status to failed
		s.repo.SMS().UpdateStatus(ctx, sms.ID.Hex(), models.StatusFailed)
		
		return common.NewServiceUnavailableError("SMS provider")
	}

	// Update status to sent
	err = s.repo.SMS().UpdateStatus(ctx, sms.ID.Hex(), models.StatusSent)
	if err != nil {
		log.Printf("Failed to update SMS status: %v", err)
	}

	log.Printf("SMS sent successfully to %s", req.PhoneNumber)
	return nil
}

// NewLogsService creates a new logs service instance
func NewLogsService(repo repository.Repository) *LogsServiceImpl {
	return &LogsServiceImpl{
		repo: repo,
	}
}

// GetLogs retrieves all OTP and callback activity logs
func (s *LogsServiceImpl) GetLogs(ctx context.Context, limit int) (map[string]interface{}, error) {
	log.Printf("Retrieving activity logs with limit: %d", limit)
	
	// Get OTP logs
	otpLogs, err := s.repo.OTP().FindAll(ctx, limit)
	if err != nil {
		log.Printf("Failed to retrieve OTP logs: %v", err)
		return nil, common.NewInternalError("Failed to retrieve OTP logs")
	}
	
	// Get callback logs
	callbackLogs, err := s.repo.Callback().FindAll(ctx, limit)
	if err != nil {
		log.Printf("Failed to retrieve callback logs: %v", err)
		return nil, common.NewInternalError("Failed to retrieve callback logs")
	}
	
	// Get SMS logs
	smsLogs, err := s.repo.SMS().FindAll(ctx, limit)
	if err != nil {
		log.Printf("Failed to retrieve SMS logs: %v", err)
		return nil, common.NewInternalError("Failed to retrieve SMS logs")
	}
	
	// Format the response
	logs := map[string]interface{}{
		"otps": map[string]interface{}{
			"count": len(otpLogs),
			"data":  otpLogs,
		},
		"callbacks": map[string]interface{}{
			"count": len(callbackLogs),
			"data":  callbackLogs,
		},
		"sms": map[string]interface{}{
			"count": len(smsLogs),
			"data":  smsLogs,
		},
		"timestamp": time.Now(),
		"total_records": len(otpLogs) + len(callbackLogs) + len(smsLogs),
	}
	
	log.Printf("Successfully retrieved logs: %d OTPs, %d callbacks, %d SMS records", 
		len(otpLogs), len(callbackLogs), len(smsLogs))
	
	return logs, nil
}

// SendOTP generates and sends a 6-digit OTP
func (s *SMSServiceImpl) SendOTP(ctx context.Context, req models.OTPRequest) (*models.OTPResponse, error) {
	log.Printf("Generating OTP for phone number: %s", req.PhoneNumber)

	// Check if OTP already exists and hasn't expired
	existingOTP, err := s.repo.OTP().FindByPhone(ctx, req.PhoneNumber)
	if err == nil && existingOTP != nil {
		// OTP exists, check if we should allow resend
		timeUntilExpiry := time.Until(existingOTP.ExpiresAt)
		if timeUntilExpiry > 2*time.Minute {
			return &models.OTPResponse{
				Success:  false,
				Message:  "OTP already sent. Please wait before requesting a new one.",
				ExpiresAt: existingOTP.ExpiresAt,
			}, nil
		}
		
		// Delete existing OTP to allow resend
		s.repo.OTP().DeleteByPhone(ctx, req.PhoneNumber)
	}

	// Generate 6-digit OTP
	otp, err := s.generateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP for %s: %v", req.PhoneNumber, err)
		return nil, common.NewInternalError("Failed to generate OTP")
	}

	// Set expiry time (5 minutes from now)
	expiry := time.Now().Add(5 * time.Minute)

	// Create OTP record
	otpRecord := &models.OTP{
		Phone:      req.PhoneNumber,
		Code:       otp,
		ExpiresAt:  expiry,
		MaxAttempts: 3,
	}

	// Store OTP in repository
	err = s.repo.OTP().Create(ctx, otpRecord)
	if err != nil {
		log.Printf("Failed to store OTP for %s: %v", req.PhoneNumber, err)
		return nil, common.NewInternalError("Failed to store OTP")
	}

	// Send OTP via SMS
	err = s.smsClient.SendOTP(ctx, req.PhoneNumber, otp)
	if err != nil {
		log.Printf("Failed to send OTP SMS to %s: %v", req.PhoneNumber, err)
		// Clean up stored OTP if SMS fails
		s.repo.OTP().DeleteByPhone(ctx, req.PhoneNumber)
		return nil, common.NewServiceUnavailableError("SMS provider")
	}

	log.Printf("OTP sent successfully to %s, expires at %v", req.PhoneNumber, expiry)

	return &models.OTPResponse{
		Success:   true,
		Message:   "OTP sent successfully",
		OTP:       otp, // In production, don't return OTP in response
		ExpiresAt: expiry,
	}, nil
}

// VerifyOTP verifies the provided OTP
func (s *SMSServiceImpl) VerifyOTP(ctx context.Context, req models.VerifyOTPRequest) (*models.VerifyOTPResponse, error) {
	log.Printf("Verifying OTP for phone number: %s", req.PhoneNumber)

	// Get stored OTP
	storedOTP, err := s.repo.OTP().FindByPhone(ctx, req.PhoneNumber)
	if err != nil || storedOTP == nil {
		log.Printf("OTP not found for %s: %v", req.PhoneNumber, err)
		return &models.VerifyOTPResponse{
			Success: false,
			Message: "OTP not found or expired. Please request a new OTP.",
			Valid:   false,
		}, nil
	}

	// Check if OTP has expired
	if time.Now().After(storedOTP.ExpiresAt) {
		log.Printf("OTP expired for %s", req.PhoneNumber)
		// Clean up expired OTP
		s.repo.OTP().DeleteByPhone(ctx, req.PhoneNumber)
		return &models.VerifyOTPResponse{
			Success: false,
			Message: "OTP expired. Please request a new OTP.",
			Valid:   false,
		}, nil
	}

	// Check if max attempts reached
	if storedOTP.Attempts >= storedOTP.MaxAttempts {
		log.Printf("Max attempts reached for %s", req.PhoneNumber)
		return &models.VerifyOTPResponse{
			Success: false,
			Message: "Maximum verification attempts reached. Please request a new OTP.",
			Valid:   false,
		}, nil
	}

	// Increment attempts
	err = s.repo.OTP().IncrementAttempts(ctx, req.PhoneNumber)
	if err != nil {
		log.Printf("Failed to increment attempts for %s: %v", req.PhoneNumber, err)
	}

	// Check if OTP matches
	if storedOTP.Code == req.OTP {
		log.Printf("OTP verified successfully for %s", req.PhoneNumber)
		
		// Delete OTP after successful verification
		s.repo.OTP().DeleteByPhone(ctx, req.PhoneNumber)
		
		return &models.VerifyOTPResponse{
			Success: true,
			Message: "OTP verified successfully",
			Valid:   true,
		}, nil
	}

	log.Printf("OTP verification failed for %s", req.PhoneNumber)
	return &models.VerifyOTPResponse{
		Success: false,
		Message: "Invalid OTP. Please try again.",
		Valid:   false,
	}, nil
}

// CleanupExpiredOTPs removes expired OTPs from storage
func (s *SMSServiceImpl) CleanupExpiredOTPs() {
	log.Println("Starting OTP cleanup routine")
	
	ctx := context.Background()
	expiredOTPs, err := s.repo.OTP().FindExpired(ctx)
	if err != nil {
		log.Printf("Failed to find expired OTPs: %v", err)
		return
	}
	
	for _, otp := range expiredOTPs {
		log.Printf("Cleaning up expired OTP for %s", otp.Phone)
		err := s.repo.OTP().DeleteByPhone(ctx, otp.Phone)
		if err != nil {
			log.Printf("Failed to delete expired OTP for %s: %v", otp.Phone, err)
		}
	}
}

// startCleanupRoutine starts the periodic cleanup of expired OTPs
func (s *SMSServiceImpl) startCleanupRoutine() {
	ticker := time.NewTicker(1 * time.Minute) // Run cleanup every minute
	defer ticker.Stop()

	for range ticker.C {
		s.CleanupExpiredOTPs()
	}
}

// generateOTP generates a random 6-digit OTP
func (s *SMSServiceImpl) generateOTP() (string, error) {
	// Generate 6 random digits
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		otp += fmt.Sprintf("%d", num.Int64())
	}
	return otp, nil
}

// NewCallbackService creates a new callback service instance
func NewCallbackService(repo repository.Repository) *CallbackServiceImpl {
	return &CallbackServiceImpl{
		repo: repo,
	}
}

// RequestCallback handles callback requests
func (s *CallbackServiceImpl) RequestCallback(ctx context.Context, req models.CallbackRequest) (*models.CallbackResponse, error) {
	log.Printf("Callback request received for phone number: %s", req.PhoneNumber)
	
	// Create callback record
	callback := &models.Callback{
		PhoneNumber: req.PhoneNumber,
		Message:     req.Message,
		Priority:    req.Priority,
		Status:      models.StatusRequested,
	}
	
	// Store callback request in database
	err := s.repo.Callback().Create(ctx, callback)
	if err != nil {
		log.Printf("Failed to store callback request for %s: %v", req.PhoneNumber, err)
		return nil, common.NewInternalError("Failed to store callback request")
	}
	
	// TODO: Placeholder for Plivo Voice API call
	// This is where you would integrate with Plivo Voice API
	// For now, just log the request
	log.Printf("Callback request logged successfully for %s. Request ID: %s", req.PhoneNumber, callback.ID.Hex())
	log.Printf("Message: %s, Priority: %s", req.Message, req.Priority)
	
	// TODO: In the future, this would make a call to Plivo Voice API
	// Example Plivo Voice API payload:
	// {
	//   "from": "+1234567890",
	//   "to": req.PhoneNumber,
	//   "answer_url": "https://your-domain.com/voice/answer",
	//   "hangup_url": "https://your-domain.com/voice/hangup",
	//   "caller_name": "SMS App"
	// }
	
	return &models.CallbackResponse{
		Success:   true,
		Message:   "Callback request received successfully",
		RequestID: callback.ID.Hex(),
		Status:    callback.Status,
		Timestamp: callback.CreatedAt,
	}, nil
}

// GetCallbackStatus retrieves the status of a callback request
func (s *CallbackServiceImpl) GetCallbackStatus(ctx context.Context, requestID string) (*models.Callback, error) {
	callback, err := s.repo.Callback().FindByID(ctx, requestID)
	if err != nil {
		return nil, common.NewNotFoundError("callback request")
	}
	return callback, nil
}

// UpdateCallbackStatus updates the status of a callback request
func (s *CallbackServiceImpl) UpdateCallbackStatus(ctx context.Context, requestID, status string) error {
	err := s.repo.Callback().UpdateStatus(ctx, requestID, status)
	if err != nil {
		return common.NewInternalError("Failed to update callback status")
	}
	return nil
} 