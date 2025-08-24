package sms_service

import (
	"context"
	"sms-app-backend/models"
)

// SMSService defines the interface for SMS operations
type SMSService interface {
	SendSMS(ctx context.Context, req models.SMSRequest) error
	SendOTP(ctx context.Context, req models.OTPRequest) (*models.OTPResponse, error)
	VerifyOTP(ctx context.Context, req models.VerifyOTPRequest) (*models.VerifyOTPResponse, error)
	CleanupExpiredOTPs()
}

// CallbackService defines the interface for callback operations
type CallbackService interface {
	RequestCallback(ctx context.Context, req models.CallbackRequest) (*models.CallbackResponse, error)
	GetCallbackStatus(ctx context.Context, requestID string) (*models.Callback, error)
	UpdateCallbackStatus(ctx context.Context, requestID, status string) error
}

// LogsService defines the interface for logs operations
type LogsService interface {
	GetLogs(ctx context.Context, limit int) (map[string]interface{}, error)
} 