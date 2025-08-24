# SMS App - Full Stack Project

A full-stack application with Next.js frontend, Golang backend, and Python AI microservice.

## 🏗️ Architecture

- **Frontend**: Next.js 14 + React + TailwindCSS
- **Backend**: Golang + Gin framework
- **AI Service**: Python + FastAPI
- **Deployment**: Vercel (Frontend) + Render/Heroku (Backend + AI)

## 📁 Project Structure

```
sms_app/
├── frontend/          # Next.js application
├── backend/           # Golang REST API
├── ai-service/        # Python AI microservice
├── docker-compose.yml # Local development
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Go 1.21+
- Python 3.11+
- Docker & Docker Compose

### Local Development

```bash
# Start all services
docker-compose up -d

# Or run individually:
cd frontend && npm run dev
cd backend && go run main.go
cd ai-service && uvicorn main:app --reload
```

### Environment Variables

Copy `.env.example` files in each directory and configure your environment variables.

## 🔧 Development

### Frontend

- Next.js 14 with App Router
- TailwindCSS for styling
- TypeScript support
- API integration with backend

### Backend

- Gin framework for REST APIs
- JWT authentication
- Database integration ready
- CORS configured

### AI Service

- FastAPI microservice
- REST endpoints for Q&A/summarization
- Async processing
- Health check endpoints

## 🚀 Deployment

### Frontend (Vercel)

```bash
cd frontend
vercel --prod
```

### Backend (Render/Heroku)

```bash
cd backend
# Configure build commands for your platform
```

### AI Service (Render/Heroku)

```bash
cd ai-service
# Configure build commands for your platform
```

## 📝 API Documentation

- Backend: `http://localhost:8080/swagger/index.html`
- AI Service: `http://localhost:8000/docs`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request
