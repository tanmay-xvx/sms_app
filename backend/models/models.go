package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Phone     string            `bson:"phone" json:"phone"`
	Email     string            `bson:"email,omitempty" json:"email,omitempty"`
	Name      string            `bson:"name,omitempty" json:"name,omitempty"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// OTP represents an OTP record
type OTP struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Phone      string            `bson:"phone" json:"phone"`
	Code       string            `bson:"code" json:"code"`
	ExpiresAt  time.Time         `bson:"expires_at" json:"expires_at"`
	Attempts   int               `bson:"attempts" json:"attempts"`
	MaxAttempts int              `bson:"max_attempts" json:"max_attempts"`
	CreatedAt  time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated_at"`
}

// SMS represents an SMS message record
type SMS struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	From        string            `bson:"from" json:"from"`
	To          string            `bson:"to" json:"to"`
	Message     string            `bson:"message" json:"message"`
	Status      string            `bson:"status" json:"status"`
	Provider    string            `bson:"provider" json:"provider"`
	ProviderID  string            `bson:"provider_id,omitempty" json:"provider_id,omitempty"`
	SentAt      time.Time         `bson:"sent_at" json:"sent_at"`
	DeliveredAt *time.Time        `bson:"delivered_at,omitempty" json:"delivered_at,omitempty"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// SMSRequest represents the request structure for sending SMS
// @Description Request structure for sending SMS
type SMSRequest struct {
	// @Description Phone number in international format (e.g., +1234567890)
	PhoneNumber string `json:"phone_number" binding:"required" example:"+1234567890"`
	// @Description SMS message content (1-160 characters)
	Message     string `json:"message" binding:"required" example:"Hello World"`
}

// OTPRequest represents the request structure for sending OTP
// @Description Request structure for sending OTP
type OTPRequest struct {
	// @Description Phone number in international format (e.g., +1234567890)
	PhoneNumber string `json:"phone_number" binding:"required" example:"+1234567890"`
}

// OTPResponse represents the response structure for OTP operations
type OTPResponse struct {
	Success   bool      `json:"success"`
	Message  string    `json:"message"`
	OTP      string    `json:"otp,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// VerifyOTPRequest represents the request structure for verifying OTP
// @Description Request structure for verifying OTP
type VerifyOTPRequest struct {
	// @Description Phone number in international format (e.g., +1234567890)
	PhoneNumber string `json:"phone_number" binding:"required" example:"+1234567890"`
	// @Description 6-digit OTP code
	OTP         string `json:"otp" binding:"required" example:"123456"`
}

// VerifyOTPResponse represents the response structure for OTP verification
type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Valid   bool   `json:"valid"`
}

// SMSResponse represents the response structure for SMS operations
type SMSResponse struct {
	Success   bool      `json:"success"`
	Message  string    `json:"message"`
	ID       string    `json:"id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// OTPStatus represents the status of an OTP
type OTPStatus struct {
	PhoneNumber string    `json:"phone_number"`
	HasActiveOTP bool     `json:"has_active_otp"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Attempts    int       `json:"attempts"`
}

// CallbackRequest represents the request structure for requesting a callback
type CallbackRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"+1234567890"`
	Message     string `json:"message,omitempty" example:"Please call me back"`
	Priority    string `json:"priority,omitempty" example:"high"`
}

// CallbackResponse represents the response structure for callback requests
type CallbackResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	RequestID string    `json:"request_id"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

// Callback represents a callback request record
type Callback struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PhoneNumber string            `bson:"phone_number" json:"phone_number"`
	Message     string            `bson:"message,omitempty" json:"message"`
	Priority    string            `bson:"priority,omitempty" json:"priority"`
	Status      string            `bson:"status" json:"status"`
	RequestedAt time.Time         `bson:"requested_at" json:"requested_at"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// PlivoCredentials represents Plivo API credentials
type PlivoCredentials struct {
	AuthID    string `json:"auth_id"`
	AuthToken string `json:"auth_token"`
	From      string `json:"from"`
}

// PlivoResponse represents the response from Plivo API
type PlivoResponse struct {
	Message     string   `json:"message"`
	MessageUUID []string `json:"message_uuid"`
	Error      string   `json:"error"`
}

// Status constants
const (
	StatusPending   = "pending"
	StatusSent      = "sent"
	StatusDelivered = "delivered"
	StatusFailed    = "failed"
	StatusRequested = "requested"
	StatusInProgress = "in_progress"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

// Provider constants
const (
	ProviderPlivo = "plivo"
	ProviderTwilio = "twilio"
) 