package mongo

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"sms-app-backend/models"
)

// MockMongoClient for testing
type MockMongoClient struct {
	otps map[string]*models.OTP
	sms  map[string]*models.SMS
	users map[string]*models.User
}

func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{
		otps:  make(map[string]*models.OTP),
		sms:   make(map[string]*models.SMS),
		users: make(map[string]*models.User),
	}
}

func (m *MockMongoClient) CreateOTP(otp *models.OTP) error {
	if otp.ID.IsZero() {
		otp.ID = primitive.NewObjectID()
	}
	m.otps[otp.ID.Hex()] = otp
	return nil
}

func (m *MockMongoClient) GetOTPByPhone(phone string) (*models.OTP, error) {
	for _, otp := range m.otps {
		if otp.Phone == phone {
			return otp, nil
		}
	}
	return nil, nil
}

func (m *MockMongoClient) UpdateOTP(otp *models.OTP) error {
	m.otps[otp.ID.Hex()] = otp
	return nil
}

func (m *MockMongoClient) DeleteOTP(id string) error {
	delete(m.otps, id)
	return nil
}

func (m *MockMongoClient) CreateSMS(sms *models.SMS) error {
	if sms.ID.IsZero() {
		sms.ID = primitive.NewObjectID()
	}
	m.sms[sms.ID.Hex()] = sms
	return nil
}

func (m *MockMongoClient) GetSMSByID(id string) (*models.SMS, error) {
	if sms, exists := m.sms[id]; exists {
		return sms, nil
	}
	return nil, nil
}

func (m *MockMongoClient) CreateUser(user *models.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	m.users[user.ID.Hex()] = user
	return nil
}

func (m *MockMongoClient) GetUserByPhone(phone string) (*models.User, error) {
	for _, user := range m.users {
		if user.Phone == phone {
			return user, nil
		}
	}
	return nil, nil
}

// Test functions
func TestOTPRepository_Create(t *testing.T) {
	mockClient := NewMockMongoClient()
	repo := &OTPRepository{}

	otp := &models.OTP{
		Phone:      "+1234567890",
		Code:       "123456",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		MaxAttempts: 3,
	}

	err := mockClient.CreateOTP(otp)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if otp.ID.IsZero() {
		t.Errorf("Expected OTP ID to be set")
	}

	if otp.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set")
	}

	if otp.UpdatedAt.IsZero() {
		t.Errorf("Expected UpdatedAt to be set")
	}
}

func TestOTPRepository_FindByPhone(t *testing.T) {
	mockClient := NewMockMongoClient()
	
	// Create test OTP
	otp := &models.OTP{
		Phone:      "+1234567890",
		Code:       "123456",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		MaxAttempts: 3,
	}
	
	err := mockClient.CreateOTP(otp)
	if err != nil {
		t.Fatalf("Failed to create test OTP: %v", err)
	}

	// Find OTP by phone
	foundOTP, err := mockClient.GetOTPByPhone("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundOTP == nil {
		t.Errorf("Expected to find OTP, got nil")
	}

	if foundOTP.Phone != "+1234567890" {
		t.Errorf("Expected phone +1234567890, got %s", foundOTP.Phone)
	}

	if foundOTP.Code != "123456" {
		t.Errorf("Expected code 123456, got %s", foundOTP.Code)
	}
}

func TestOTPRepository_Update(t *testing.T) {
	mockClient := NewMockMongoClient()
	
	// Create test OTP
	otp := &models.OTP{
		Phone:      "+1234567890",
		Code:       "123456",
		ExpiresAt:  time.Now().Add(5 * time.Minute),
		MaxAttempts: 3,
	}
	
	err := mockClient.CreateOTP(otp)
	if err != nil {
		t.Fatalf("Failed to create test OTP: %v", err)
	}

	// Update OTP
	otp.Code = "654321"
	otp.Attempts = 1
	
	err = mockClient.UpdateOTP(otp)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	updatedOTP, err := mockClient.GetOTPByPhone("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if updatedOTP.Code != "654321" {
		t.Errorf("Expected updated code 654321, got %s", updatedOTP.Code)
	}

	if updatedOTP.Attempts != 1 {
		t.Errorf("Expected attempts 1, got %d", updatedOTP.Attempts)
	}
}

func TestSMSRepository_Create(t *testing.T) {
	mockClient := NewMockMongoClient()
	
	sms := &models.SMS{
		From:     "+0987654321",
		To:       "+1234567890",
		Message:  "Test message",
		Status:   models.StatusPending,
		Provider: models.ProviderPlivo,
	}

	err := mockClient.CreateSMS(sms)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if sms.ID.IsZero() {
		t.Errorf("Expected SMS ID to be set")
	}

	if sms.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set")
	}

	if sms.SentAt.IsZero() {
		t.Errorf("Expected SentAt to be set")
	}
}

func TestUserRepository_Create(t *testing.T) {
	mockClient := NewMockMongoClient()
	
	user := &models.User{
		Phone: "+1234567890",
		Email: "test@example.com",
		Name:  "Test User",
	}

	err := mockClient.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID.IsZero() {
		t.Errorf("Expected User ID to be set")
	}

	if user.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set")
	}
}

func TestUserRepository_FindByPhone(t *testing.T) {
	mockClient := NewMockMongoClient()
	
	// Create test user
	user := &models.User{
		Phone: "+1234567890",
		Email: "test@example.com",
		Name:  "Test User",
	}
	
	err := mockClient.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Find user by phone
	foundUser, err := mockClient.GetUserByPhone("+1234567890")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundUser == nil {
		t.Errorf("Expected to find user, got nil")
	}

	if foundUser.Phone != "+1234567890" {
		t.Errorf("Expected phone +1234567890, got %s", foundUser.Phone)
	}

	if foundUser.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", foundUser.Email)
	}
} 