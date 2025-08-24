# 🚀 SMS Application Deployment Guide

This guide covers deployment of all three services to production environments.

## 📋 **Prerequisites**

- [Vercel](https://vercel.com) account (for frontend)
- [Render](https://render.com) account (for backend & AI service)
- [MongoDB Atlas](https://mongodb.com/atlas) account (or local MongoDB)
- [Plivo](https://plivo.com) account for SMS API
- [OpenAI](https://openai.com) API key

## 🎯 **Service Architecture**

```
Frontend (Vercel) → Backend (Render) → MongoDB (Local/Atlas)
                ↓
            AI Service (Render) → OpenAI API
```

---

## 🌐 **1. Frontend → Vercel**

### **Step 1: Prepare Frontend**

```bash
cd frontend
npm run build  # Test build locally
```

### **Step 2: Deploy to Vercel**

```bash
# Install Vercel CLI
npm i -g vercel

# Login to Vercel
vercel login

# Deploy
vercel --prod
```

### **Step 3: Configure Environment Variables**

In Vercel Dashboard → Project Settings → Environment Variables:

| Variable                     | Value                                  |
| ---------------------------- | -------------------------------------- |
| `NEXT_PUBLIC_API_URL`        | `https://your-backend.onrender.com`    |
| `NEXT_PUBLIC_AI_SERVICE_URL` | `https://your-ai-service.onrender.com` |

### **Step 4: Update CORS Origins**

Add your Vercel domain to backend and AI service CORS configurations.

---

## 🔧 **2. Golang Backend → Render**

### **Step 1: Prepare Backend**

```bash
cd backend
go build -o sms-backend main.go  # Test build
```

### **Step 2: Create Render Service**

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Configure service:

**Build Command:**

```bash
go build -o sms-backend main.go
```

**Start Command:**

```bash
./sms-backend
```

**Environment:**

```bash
go 1.21
```

### **Step 3: Configure Environment Variables**

In Render Dashboard → Environment:

| Variable                  | Value                                     |
| ------------------------- | ----------------------------------------- |
| `ENVIRONMENT`             | `production`                              |
| `GIN_MODE`                | `release`                                 |
| `PORT`                    | `10000`                                   |
| `AI_SERVICE_URL`          | `https://your-ai-service.onrender.com`    |
| `JWT_SECRET`              | `your_secure_jwt_secret_here`             |
| `CORS_ORIGIN`             | `https://your-frontend.vercel.app`        |
| `ADDITIONAL_CORS_ORIGINS` | `https://your-ai-service.onrender.com`    |
| `MONGODB_URI`             | `mongodb://your-mongo-host:27017/sms_app` |
| `PLIVO_AUTH_ID`           | `your_production_plivo_auth_id`           |
| `PLIVO_AUTH_TOKEN`        | `your_production_plivo_auth_token`        |
| `PLIVO_FROM_NUMBER`       | `your_production_plivo_phone_number`      |

### **Step 4: Health Check**

Render will automatically use the `/health` endpoint for health checks.

---

## 🤖 **3. Python AI Service → Render**

### **Step 1: Prepare AI Service**

```bash
cd ai-service
pip install -r requirements.txt
python -m uvicorn main:app --host 0.0.0.0 --port 8000  # Test locally
```

### **Step 2: Create Render Service**

1. Go to [Render Dashboard](https://dashboard.render.com)
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Configure service:

**Build Command:**

```bash
pip install -r requirements.txt
```

**Start Command:**

```bash
uvicorn main:app --host 0.0.0.0 --port $PORT
```

**Environment:**

```bash
python 3.11
```

### **Step 3: Configure Environment Variables**

In Render Dashboard → Environment:

| Variable         | Value                                                                |
| ---------------- | -------------------------------------------------------------------- |
| `ENVIRONMENT`    | `production`                                                         |
| `PORT`           | `10000`                                                              |
| `OPENAI_API_KEY` | `your_production_openai_api_key`                                     |
| `MODEL_NAME`     | `gpt-4`                                                              |
| `CORS_ORIGINS`   | `https://your-frontend.vercel.app,https://your-backend.onrender.com` |

### **Step 4: Health Check**

Render will automatically use the `/health` endpoint for health checks.

---

## 🗄️ **4. MongoDB Configuration**

### **Option A: Local MongoDB (Exposed)**

```bash
# Install MongoDB locally
brew install mongodb/brew/mongodb-community

# Start MongoDB
brew services start mongodb/brew/mongodb-community

# Expose MongoDB to network (for Render access)
# Edit /usr/local/etc/mongod.conf
net:
  port: 27017
  bindIp: 0.0.0.0  # Allow external connections

# Restart MongoDB
brew services restart mongodb/brew/mongodb-community
```

**Connection String:**

```
mongodb://your-local-ip:27017/sms_app
```

### **Option B: MongoDB Atlas (Recommended)**

1. Create cluster in [MongoDB Atlas](https://mongodb.com/atlas)
2. Create database user
3. Get connection string:

```
mongodb+srv://username:password@cluster.mongodb.net/sms_app
```

---

## 🔐 **5. Security Configuration**

### **Generate Secure Secrets**

```bash
# JWT Secret (64 characters)
openssl rand -base64 64

# Environment-specific secrets
openssl rand -base64 32
```

### **Environment Variable Security**

- ✅ **Vercel**: Environment variables are encrypted
- ✅ **Render**: Environment variables are encrypted
- ❌ **Never commit** `.env` files to Git
- ✅ **Use** `.env.example` for documentation

### **CORS Security**

- ✅ **Production**: Only allow specific domains
- ❌ **Development**: Can allow localhost
- ✅ **Validate** all incoming origins

---

## 🧪 **6. Testing Deployment**

### **Health Checks**

```bash
# Backend
curl https://your-backend.onrender.com/health

# AI Service
curl https://your-ai-service.onrender.com/health

# Frontend
curl https://your-frontend.vercel.app
```

### **API Tests**

```bash
# Test OTP sending
curl -X POST https://your-backend.onrender.com/api/sms/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'

# Test AI chat
curl -X POST https://your-ai-service.onrender.com/chat \
  -H "Content-Type: application/json" \
  -d '{"question": "Hello"}'
```

---

## 📊 **7. Monitoring & Logs**

### **Render Logs**

- View real-time logs in Render Dashboard
- Set up log forwarding to external services
- Monitor service health and performance

### **Vercel Analytics**

- Enable Vercel Analytics for frontend monitoring
- Track performance metrics
- Monitor user experience

### **MongoDB Monitoring**

- Use MongoDB Compass for database monitoring
- Set up alerts for connection issues
- Monitor query performance

---

## 🚨 **8. Troubleshooting**

### **Common Issues**

**CORS Errors:**

```bash
# Check CORS configuration
# Verify allowed origins in both services
# Test with Postman/curl first
```

**Database Connection:**

```bash
# Verify MongoDB URI
# Check network access
# Test connection locally
```

**Build Failures:**

```bash
# Check Go/Python versions
# Verify dependencies
# Test builds locally first
```

### **Debug Commands**

```bash
# Check service status
curl -v https://your-service.onrender.com/health

# Test database connection
mongosh "your-mongodb-uri"

# Verify environment variables
echo $VARIABLE_NAME
```

---

## 🔄 **9. CI/CD Setup**

### **Automatic Deployments**

- Connect GitHub repository to Render
- Enable automatic deployments on push to main
- Set up preview deployments for pull requests

### **Environment Promotion**

- Use different environments (staging/production)
- Test in staging before production
- Use feature flags for gradual rollouts

---

## 📈 **10. Performance Optimization**

### **Backend (Golang)**

- Enable GIN release mode
- Use connection pooling for MongoDB
- Implement rate limiting
- Add caching layers

### **AI Service (Python)**

- Use async/await for I/O operations
- Implement request queuing
- Add response caching
- Monitor OpenAI API usage

### **Frontend (Next.js)**

- Enable static generation where possible
- Implement lazy loading
- Use CDN for static assets
- Optimize bundle size

---

## 🎯 **Deployment Checklist**

- [ ] All services build successfully locally
- [ ] Environment variables configured in production
- [ ] CORS origins updated for production domains
- [ ] Database connection tested
- [ ] Health check endpoints working
- [ ] API endpoints tested
- [ ] Frontend-backend integration verified
- [ ] AI service integration tested
- [ ] Monitoring and logging configured
- [ ] Security measures implemented
- [ ] Performance tested under load

---

## 📞 **Support**

For deployment issues:

1. Check service logs in Render/Vercel dashboards
2. Verify environment variable configuration
3. Test endpoints individually
4. Check CORS and network connectivity
5. Review security group and firewall settings

**Happy Deploying! 🚀**
