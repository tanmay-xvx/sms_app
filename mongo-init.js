// MongoDB initialization script for SMS App
// This script runs when the MongoDB container starts for the first time

// Create database and user
db = db.getSiblingDB('sms_app');

// Create collections with proper indexes
db.createCollection('otps');
db.createCollection('sms_messages');
db.createCollection('users');
db.createCollection('callbacks');

// Create indexes for better performance
db.otps.createIndex({ "phone_number": 1 });
db.otps.createIndex({ "created_at": 1 });
db.otps.createIndex({ "expires_at": 1 }, { expireAfterSeconds: 0 });

db.sms_messages.createIndex({ "phone_number": 1 });
db.sms_messages.createIndex({ "created_at": 1 });

db.users.createIndex({ "phone_number": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { sparse: true });

db.callbacks.createIndex({ "phone_number": 1 });
db.callbacks.createIndex({ "status": 1 });
db.callbacks.createIndex({ "requested_at": 1 });

// Insert sample data for development
db.otps.insertOne({
    phone_number: "+1234567890",
    otp: "123456",
    created_at: new Date(),
    expires_at: new Date(Date.now() + 5 * 60 * 1000), // 5 minutes
    attempts: 0,
    verified: false
});

db.sms_messages.insertOne({
    phone_number: "+1234567890",
    message: "Welcome to SMS App!",
    status: "sent",
    created_at: new Date(),
    message_type: "welcome"
});

db.callbacks.insertOne({
    phone_number: "+1234567890",
    message: "Sample callback request",
    priority: "normal",
    status: "requested",
    requested_at: new Date(),
    created_at: new Date()
});

print("MongoDB initialized successfully for SMS App!");
print("Database: sms_app");
print("Collections: otps, sms_messages, users, callbacks");
print("Indexes created for optimal performance"); 