#!/bin/bash

echo "🚀 Setting up SMS App Full-Stack Project"
echo "========================================"

# Check if required tools are installed
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    fi
}

echo "🔍 Checking prerequisites..."
check_command "node"
check_command "npm"
check_command "go"
check_command "python3"
check_command "pip3"
check_command "docker"
check_command "docker-compose"

echo "✅ All prerequisites are installed!"

# Create .env files from examples
echo "📝 Setting up environment files..."

if [ ! -f "frontend/.env" ]; then
    cp frontend/env.example frontend/.env
    echo "✅ Created frontend/.env"
fi

if [ ! -f "backend/.env" ]; then
    cp backend/env.example backend/.env
    echo "✅ Created backend/.env"
fi

if [ ! -f "ai-service/.env" ]; then
    cp ai-service/env.example ai-service/.env
    echo "✅ Created ai-service/.env"
fi

# Install dependencies
echo "📦 Installing dependencies..."

echo "Installing frontend dependencies..."
cd frontend && npm install && cd ..

echo "Installing backend dependencies..."
cd backend && go mod tidy && cd ..

echo "Installing AI service dependencies..."
cd ai-service && pip3 install -r requirements.txt && cd ..

echo "✅ Dependencies installed!"

# Build Docker images
echo "🐳 Building Docker images..."
docker-compose build

echo "✅ Docker images built!"

echo ""
echo "🎉 Setup complete! Here's what you can do next:"
echo ""
echo "1. Start all services:"
echo "   make dev"
echo "   # or"
echo "   docker-compose up -d"
echo ""
echo "2. Start individual services:"
echo "   make dev-frontend    # Frontend only"
echo "   make dev-backend     # Backend only"
echo "   make dev-ai          # AI service only"
echo ""
echo "3. View service status:"
echo "   make status"
echo ""
echo "4. View logs:"
echo "   make logs"
echo ""
echo "5. Stop services:"
echo "   make docker-down"
echo ""
echo "📚 Documentation:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   Backend Docs: http://localhost:8080/swagger/index.html"
echo "   AI Service: http://localhost:8000"
echo "   AI Service Docs: http://localhost:8000/docs"
echo ""
echo "⚠️  Don't forget to:"
echo "   - Update .env files with your actual API keys"
echo "   - Configure CORS origins for production"
echo "   - Set up proper JWT secrets"
echo ""
echo "Happy coding! 🚀" 