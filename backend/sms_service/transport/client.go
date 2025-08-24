package transport

import (
	"context"
	"sms-app-backend/models"
)

// SMSClient defines the interface for SMS service clients
type SMSClient interface {
	SendSMS(ctx context.Context, to, message string) error
	SendOTP(ctx context.Context, to, otp string) error
	GetProvider() string
}

// PlivoClient implements SMSClient for Plivo SMS service
type PlivoClient struct {
	authID    string
	authToken string
	from      string
	baseURL   string
}

// NewPlivoClient creates a new Plivo client
func NewPlivoClient(authID, authToken, from string) *PlivoClient {
	return &PlivoClient{
		authID:    authID,
		authToken: authToken,
		from:      from,
		baseURL:   "https://api.plivo.com/v1/Account/" + authID + "/Message/",
	}
}

// SendSMS sends an SMS message via Plivo
func (pc *PlivoClient) SendSMS(ctx context.Context, to, message string) error {
	// Implementation would use HTTP client to call Plivo API
	// For now, return nil to indicate success
	return nil
}

// SendOTP sends an OTP message via Plivo
func (pc *PlivoClient) SendOTP(ctx context.Context, to, otp string) error {
	message := "Your OTP is: " + otp + ". Valid for 5 minutes. Do not share this code."
	return pc.SendSMS(ctx, to, message)
}

// GetProvider returns the provider name
func (pc *PlivoClient) GetProvider() string {
	return models.ProviderPlivo
}

// MockClient implements SMSClient for testing
type MockClient struct {
	provider string
}

// NewMockClient creates a new mock SMS client
func NewMockClient(provider string) *MockClient {
	return &MockClient{provider: provider}
}

// SendSMS mock implementation
func (mc *MockClient) SendSMS(ctx context.Context, to, message string) error {
	return nil
}

// SendOTP mock implementation
func (mc *MockClient) SendOTP(ctx context.Context, to, otp string) error {
	return nil
}

// GetProvider returns the provider name
func (mc *MockClient) GetProvider() string {
	return mc.provider
} 