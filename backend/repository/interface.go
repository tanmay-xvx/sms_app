package repository

import (
	"context"
	"time"

	"sms-app-backend/models"
)

// OTPRepository defines the interface for OTP storage operations
type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTP) error
	FindByPhone(ctx context.Context, phone string) (*models.OTP, error)
	Update(ctx context.Context, otp *models.OTP) error
	Delete(ctx context.Context, id string) error
	DeleteByPhone(ctx context.Context, phone string) error
	FindExpired(ctx context.Context) ([]*models.OTP, error)
	IncrementAttempts(ctx context.Context, phone string) error
	FindAll(ctx context.Context, limit int) ([]*models.OTP, error)
}

// SMSRepository defines the interface for SMS storage operations
type SMSRepository interface {
	Create(ctx context.Context, sms *models.SMS) error
	FindByID(ctx context.Context, id string) (*models.SMS, error)
	FindByPhone(ctx context.Context, phone string, limit int) ([]*models.SMS, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdateDeliveryTime(ctx context.Context, id string, deliveredAt time.Time) error
	FindByStatus(ctx context.Context, status string, limit int) ([]*models.SMS, error)
	FindAll(ctx context.Context, limit int) ([]*models.SMS, error)
}

// UserRepository defines the interface for user storage operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

// CallbackRepository defines the interface for callback storage operations
type CallbackRepository interface {
	Create(ctx context.Context, callback *models.Callback) error
	FindByID(ctx context.Context, id string) (*models.Callback, error)
	FindByPhone(ctx context.Context, phone string, limit int) ([]*models.Callback, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	FindByStatus(ctx context.Context, status string, limit int) ([]*models.Callback, error)
	FindAll(ctx context.Context, limit int) ([]*models.Callback, error)
}

// Repository defines the main repository interface
type Repository interface {
	OTP() OTPRepository
	SMS() SMSRepository
	User() UserRepository
	Callback() CallbackRepository
	Close() error
} 