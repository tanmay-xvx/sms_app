#!/bin/bash

# 🚀 SMS Application Production Deployment Script
# This script helps deploy all three services to production

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
FRONTEND_DIR="frontend"
BACKEND_DIR="backend"
AI_SERVICE_DIR="ai-service"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check Node.js
    if ! command_exists node; then
        print_error "Node.js is not installed. Please install Node.js 18+"
        exit 1
    fi
    
    # Check npm
    if ! command_exists npm; then
        print_error "npm is not installed. Please install npm"
        exit 1
    fi
    
    # Check Go
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21+"
        exit 1
    fi
    
    # Check Python
    if ! command_exists python3; then
        print_error "Python 3 is not installed. Please install Python 3.11+"
        exit 1
    fi
    
    # Check Vercel CLI
    if ! command_exists vercel; then
        print_warning "Vercel CLI not found. Installing..."
        npm install -g vercel
    fi
    
    print_success "All prerequisites are satisfied!"
}

# Function to build frontend
build_frontend() {
    print_status "Building frontend..."
    cd "$FRONTEND_DIR"
    
    # Install dependencies
    npm install
    
    # Build the application
    npm run build
    
    print_success "Frontend built successfully!"
    cd ..
}

# Function to build backend
build_backend() {
    print_status "Building backend..."
    cd "$BACKEND_DIR"
    
    # Download dependencies
    go mod download
    
    # Build the application
    go build -o sms-backend main.go
    
    print_success "Backend built successfully!"
    cd ..
}

# Function to build AI service
build_ai_service() {
    print_status "Building AI service..."
    cd "$AI_SERVICE_DIR"
    
    # Install dependencies
    pip3 install -r requirements.txt
    
    # Test the service
    python3 -c "import fastapi; print('FastAPI imported successfully')"
    
    print_success "AI service built successfully!"
    cd ..
}

# Function to deploy frontend to Vercel
deploy_frontend() {
    print_status "Deploying frontend to Vercel..."
    cd "$FRONTEND_DIR"
    
    # Check if already logged in
    if ! vercel whoami >/dev/null 2>&1; then
        print_warning "Please login to Vercel first:"
        vercel login
    fi
    
    # Deploy to production
    vercel --prod --yes
    
    print_success "Frontend deployed to Vercel!"
    cd ..
}

# Function to create Render deployment files
create_render_files() {
    print_status "Creating Render deployment files..."
    
    # Create backend render.yaml
    cat > "$BACKEND_DIR/render.yaml" << EOF
services:
  - type: web
    name: sms-backend
    env: go
    buildCommand: go build -o sms-backend main.go
    startCommand: ./sms-backend
    envVars:
      - key: ENVIRONMENT
        value: production
      - key: GIN_MODE
        value: release
      - key: PORT
        value: 10000
      - key: JWT_SECRET
        generateValue: true
      - key: CORS_ORIGIN
        sync: false
      - key: MONGODB_URI
        sync: false
      - key: PLIVO_AUTH_ID
        sync: false
      - key: PLIVO_AUTH_TOKEN
        sync: false
      - key: PLIVO_FROM_NUMBER
        sync: false
EOF

    # Create AI service render.yaml
    cat > "$AI_SERVICE_DIR/render.yaml" << EOF
services:
  - type: web
    name: sms-ai-service
    env: python
    buildCommand: pip install -r requirements.txt
    startCommand: uvicorn main:app --host 0.0.0.0 --port \$PORT
    envVars:
      - key: ENVIRONMENT
        value: production
      - key: PORT
        value: 10000
      - key: OPENAI_API_KEY
        sync: false
      - key: MODEL_NAME
        value: gpt-4
      - key: CORS_ORIGINS
        sync: false
EOF

    print_success "Render deployment files created!"
}

# Function to show deployment instructions
show_instructions() {
    echo
    print_status "🎯 Deployment Instructions:"
    echo
    echo "1. 🌐 Frontend (Vercel):"
    echo "   - Frontend has been deployed to Vercel"
    echo "   - Set environment variables in Vercel dashboard:"
    echo "     • NEXT_PUBLIC_API_URL=https://your-backend.onrender.com"
    echo "     • NEXT_PUBLIC_AI_SERVICE_URL=https://your-ai-service.onrender.com"
    echo
    echo "2. 🔧 Backend (Render):"
    echo "   - Go to https://dashboard.render.com"
    echo "   - Create new Web Service"
    echo "   - Connect your GitHub repository"
    echo "   - Use render.yaml from backend/ directory"
    echo "   - Set environment variables in Render dashboard"
    echo
    echo "3. 🤖 AI Service (Render):"
    echo "   - Go to https://dashboard.render.com"
    echo "   - Create new Web Service"
    echo "   - Connect your GitHub repository"
    echo "   - Use render.yaml from ai-service/ directory"
    echo "   - Set environment variables in Render dashboard"
    echo
    echo "4. 🗄️ MongoDB:"
    echo "   - Use MongoDB Atlas (recommended) or expose local MongoDB"
    echo "   - Update MONGODB_URI in backend environment variables"
    echo
    echo "5. 🔐 Security:"
    echo "   - Generate secure JWT_SECRET: openssl rand -base64 64"
    echo "   - Set production Plivo credentials"
    echo   "   - Set production OpenAI API key"
    echo
}

# Main deployment flow
main() {
    echo "🚀 SMS Application Production Deployment"
    echo "========================================"
    echo
    
    # Check prerequisites
    check_prerequisites
    
    # Build all services
    build_frontend
    build_backend
    build_ai_service
    
    # Deploy frontend
    deploy_frontend
    
    # Create Render deployment files
    create_render_files
    
    # Show instructions
    show_instructions
    
    print_success "Deployment preparation completed!"
    print_status "Follow the instructions above to complete deployment to Render"
}

# Run main function
main "$@" 