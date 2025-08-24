package sms_service

import (
	"context"
	"testing"
	"time"
)

// MockPlivoClient for testing
type MockPlivoClient struct{}

func (m *MockPlivoClient) SendSMS(to, message string) error {
	return nil
}

func (m *MockPlivoClient) SendOTP(to, otp string) error {
	return nil
}

func TestSendOTP(t *testing.T) {
	// Create mock components
	otpRepo := NewInMemoryOTPRepository()
	mockPlivo := &MockPlivoClient{}
	
	// Create service
	service := NewSMSService(otpRepo, mockPlivo)
	
	// Test OTP generation
	req := OTPRequest{PhoneNumber: "+1234567890"}
	response, err := service.SendOTP(context.Background(), req)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !response.Success {
		t.Errorf("Expected success=true, got %v", response.Success)
	}
	
	if response.OTP == "" {
		t.Errorf("Expected OTP to be generated, got empty string")
	}
	
	if len(response.OTP) != 6 {
		t.Errorf("Expected 6-digit OTP, got %d digits", len(response.OTP))
	}
}

func TestOTPExpiry(t *testing.T) {
	otpRepo := NewInMemoryOTPRepository()
	mockPlivo := &MockPlivoClient{}
	service := NewSMSService(otpRepo, mockPlivo)
	
	// Send OTP
	req := OTPRequest{PhoneNumber: "+1234567890"}
	response, err := service.SendOTP(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to send OTP: %v", err)
	}
	
	// Verify OTP is stored
	otp, expiry, err := otpRepo.GetOTP("+1234567890")
	if err != nil {
		t.Errorf("Expected OTP to be stored, got error: %v", err)
	}
	
	if otp != response.OTP {
		t.Errorf("Expected stored OTP to match generated OTP")
	}
	
	// Check expiry is set to 5 minutes from now
	expectedExpiry := time.Now().Add(5 * time.Minute)
	if time.Until(expectedExpiry) > 10*time.Second {
		t.Errorf("Expected expiry to be approximately 5 minutes from now")
	}
}

func TestVerifyOTP(t *testing.T) {
	otpRepo := NewInMemoryOTPRepository()
	mockPlivo := &MockPlivoClient{}
	service := NewSMSService(otpRepo, mockPlivo)
	
	// Send OTP first
	req := OTPRequest{PhoneNumber: "+1234567890"}
	response, err := service.SendOTP(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to send OTP: %v", err)
	}
	
	// Verify with correct OTP
	verifyReq := VerifyOTPRequest{
		PhoneNumber: "+1234567890",
		OTP:         response.OTP,
	}
	
	verifyResp, err := service.VerifyOTP(context.Background(), verifyReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if !verifyResp.Success {
		t.Errorf("Expected verification to succeed, got %v", verifyResp.Success)
	}
	
	if !verifyResp.Valid {
		t.Errorf("Expected OTP to be valid, got %v", verifyResp.Valid)
	}
}

func TestInvalidOTP(t *testing.T) {
	otpRepo := NewInMemoryOTPRepository()
	mockPlivo := &MockPlivoClient{}
	service := NewSMSService(otpRepo, mockPlivo)
	
	// Send OTP first
	req := OTPRequest{PhoneNumber: "+1234567890"}
	_, err := service.SendOTP(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to send OTP: %v", err)
	}
	
	// Verify with incorrect OTP
	verifyReq := VerifyOTPRequest{
		PhoneNumber: "+1234567890",
		OTP:         "000000",
	}
	
	verifyResp, err := service.VerifyOTP(context.Background(), verifyReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if verifyResp.Success {
		t.Errorf("Expected verification to fail, got %v", verifyResp.Success)
	}
	
	if verifyResp.Valid {
		t.Errorf("Expected OTP to be invalid, got %v", verifyResp.Valid)
	}
} 