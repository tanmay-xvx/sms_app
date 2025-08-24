package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sms-app-backend/models"
	"sms-app-backend/repository"
)

// Repository implements the repository.Repository interface
type Repository struct {
	client       *mongo.Client
	database     *mongo.Database
	otpRepo      *OTPRepository
	smsRepo      *SMSRepository
	userRepo     *UserRepository
	callbackRepo *CallbackRepository
}

// NewRepository creates a new MongoDB repository
func NewRepository(uri, dbName string) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := client.Database(dbName)

	repo := &Repository{
		client:   client,
		database: database,
	}

	// Initialize sub-repositories
	repo.otpRepo = NewOTPRepository(database)
	repo.smsRepo = NewSMSRepository(database)
	repo.userRepo = NewUserRepository(database)
	repo.callbackRepo = NewCallbackRepository(database)

	return repo, nil
}

// OTP returns the OTP repository
func (r *Repository) OTP() repository.OTPRepository {
	return r.otpRepo
}

// SMS returns the SMS repository
func (r *Repository) SMS() repository.SMSRepository {
	return r.smsRepo
}

// User returns the user repository
func (r *Repository) User() repository.UserRepository {
	return r.userRepo
}

// Callback returns the callback repository
func (r *Repository) Callback() repository.CallbackRepository {
	return r.callbackRepo
}

// Close closes the MongoDB connection
func (r *Repository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}

// OTPRepository implements repository.OTPRepository
type OTPRepository struct {
	collection *mongo.Collection
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *mongo.Database) *OTPRepository {
	collection := db.Collection("otps")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Index on phone number
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		// Index might already exist
	}

	// Index on expiry for cleanup
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		// Index might already exist
	}

	return &OTPRepository{collection: collection}
}

// Create stores a new OTP
func (r *OTPRepository) Create(ctx context.Context, otp *models.OTP) error {
	otp.CreatedAt = time.Now()
	otp.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, otp)
	if err != nil {
		return err
	}
	
	otp.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByPhone finds an OTP by phone number
func (r *OTPRepository) FindByPhone(ctx context.Context, phone string) (*models.OTP, error) {
	var otp models.OTP
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// Update updates an existing OTP
func (r *OTPRepository) Update(ctx context.Context, otp *models.OTP) error {
	otp.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": otp.ID},
		bson.M{"$set": otp},
	)
	return err
}

// Delete deletes an OTP by ID
func (r *OTPRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

// CallbackRepository implements repository.CallbackRepository
type CallbackRepository struct {
	collection *mongo.Collection
}

// NewCallbackRepository creates a new callback repository
func NewCallbackRepository(db *mongo.Database) *CallbackRepository {
	collection := db.Collection("callbacks")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Index on phone number
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "phone_number", Value: 1}},
	})
	if err != nil {
		// Index might already exist
	}

	// Index on status
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	})
	if err != nil {
		// Index might already exist
	}

	// Index on requested_at for sorting
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "requested_at", Value: -1}},
	})
	if err != nil {
		// Index might already exist
	}

	return &CallbackRepository{collection: collection}
}

// Create stores a new callback request
func (r *CallbackRepository) Create(ctx context.Context, callback *models.Callback) error {
	callback.CreatedAt = time.Now()
	callback.UpdatedAt = time.Now()
	callback.RequestedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, callback)
	if err != nil {
		return err
	}
	
	callback.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a callback by ID
func (r *CallbackRepository) FindByID(ctx context.Context, id string) (*models.Callback, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var callback models.Callback
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&callback)
	if err != nil {
		return nil, err
	}
	return &callback, nil
}

// FindByPhone finds callback requests by phone number
func (r *CallbackRepository) FindByPhone(ctx context.Context, phone string, limit int) ([]*models.Callback, error) {
	opts := options.Find().SetSort(bson.D{{Key: "requested_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{"phone_number": phone}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var callbacks []*models.Callback
	if err = cursor.All(ctx, &callbacks); err != nil {
		return nil, err
	}
	
	return callbacks, nil
}

// UpdateStatus updates the status of a callback
func (r *CallbackRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}

// FindByStatus finds callback requests by status
func (r *CallbackRepository) FindByStatus(ctx context.Context, status string, limit int) ([]*models.Callback, error) {
	opts := options.Find().SetSort(bson.D{{Key: "requested_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	
	defer cursor.Close(ctx)

	var callbacks []*models.Callback
	if err = cursor.All(ctx, &callbacks); err != nil {
		return nil, err
	}
	
	return callbacks, nil
}

// FindAll finds all callback requests with a limit
func (r *CallbackRepository) FindAll(ctx context.Context, limit int) ([]*models.Callback, error) {
	opts := options.Find().SetSort(bson.D{{Key: "requested_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var callbacks []*models.Callback
	if err = cursor.All(ctx, &callbacks); err != nil {
		return nil, err
	}
	return callbacks, nil
}

// DeleteByPhone deletes an OTP by phone number
func (r *OTPRepository) DeleteByPhone(ctx context.Context, phone string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"phone": phone})
	return err
}

// FindExpired finds all expired OTPs
func (r *OTPRepository) FindExpired(ctx context.Context) ([]*models.OTP, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"expires_at": bson.M{"$lt": time.Now()}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var otps []*models.OTP
	if err = cursor.All(ctx, &otps); err != nil {
		return nil, err
	}
	
	return otps, nil
}

// FindAll finds all OTPs with a limit
func (r *OTPRepository) FindAll(ctx context.Context, limit int) ([]*models.OTP, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var otps []*models.OTP
	if err = cursor.All(ctx, &otps); err != nil {
		return nil, err
	}
	return otps, nil
}

// IncrementAttempts increments the attempt counter for a phone number
func (r *OTPRepository) IncrementAttempts(ctx context.Context, phone string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"phone": phone},
		bson.M{"$inc": bson.M{"attempts": 1}, "$set": bson.M{"updated_at": time.Now()}},
	)
	return err
}

// SMSRepository implements repository.SMSRepository
type SMSRepository struct {
	collection *mongo.Collection
}

// NewSMSRepository creates a new SMS repository
func NewSMSRepository(db *mongo.Database) *SMSRepository {
	collection := db.Collection("sms")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Index on phone numbers
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "to", Value: 1}},
	})
	if err != nil {
		// Index might already exist
	}

	// Index on status
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	})
	if err != nil {
		// Index might already exist
	}

	return &SMSRepository{collection: collection}
}

// Create stores a new SMS
func (r *SMSRepository) Create(ctx context.Context, sms *models.SMS) error {
	sms.CreatedAt = time.Now()
	sms.UpdatedAt = time.Now()
	sms.SentAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, sms)
	if err != nil {
		return err
	}
	
	sms.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds an SMS by ID
func (r *SMSRepository) FindByID(ctx context.Context, id string) (*models.SMS, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var sms models.SMS
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&sms)
	if err != nil {
		return nil, err
	}
	return &sms, nil
}

// FindByPhone finds SMS messages by phone number
func (r *SMSRepository) FindByPhone(ctx context.Context, phone string, limit int) ([]*models.SMS, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{"to": phone}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sms []*models.SMS
	if err = cursor.All(ctx, &sms); err != nil {
		return nil, err
	}
	
	return sms, nil
}

// UpdateStatus updates the status of an SMS
func (r *SMSRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}

// UpdateDeliveryTime updates the delivery time of an SMS
func (r *SMSRepository) UpdateDeliveryTime(ctx context.Context, id string, deliveredAt time.Time) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"delivered_at": deliveredAt, "updated_at": time.Now()}},
	)
	return err
}

// FindByStatus finds SMS messages by status
func (r *SMSRepository) FindByStatus(ctx context.Context, status string, limit int) ([]*models.SMS, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sms []*models.SMS
	if err = cursor.All(ctx, &sms); err != nil {
		return nil, err
	}
	
	return sms, nil
}

// FindAll finds all SMS messages with a limit
func (r *SMSRepository) FindAll(ctx context.Context, limit int) ([]*models.SMS, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(int64(limit))
	
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sms []*models.SMS
	if err = cursor.All(ctx, &sms); err != nil {
		return nil, err
	}
	return sms, nil
}

// UserRepository implements repository.UserRepository
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection("users")
	
	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// Index on phone number
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		// Index might already exist
	}

	// Index on email
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "email", Value: 1}},
	})
	if err != nil {
		// Index might already exist
	}

	return &UserRepository{collection: collection}
}

// Create stores a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone finds a user by phone number
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"phone": phone}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
} 