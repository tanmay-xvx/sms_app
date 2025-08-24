# ✅ Production Deployment Checklist

## 🎯 **Pre-Deployment Verification**

### **Local Testing** ✅

- [x] Frontend builds successfully: `npm run build`
- [x] Backend builds successfully: `go build -o sms-backend main.go`
- [x] AI service dependencies installed: `pip install -r requirements.txt`
- [x] All services running locally
- [x] Health check endpoints working
- [x] Complete flow tested (OTP → Verify → Dashboard → AI Chat)

---

## 🌐 **1. Frontend → Vercel**

### **Build & Deploy** ✅

- [ ] Run: `cd frontend && npm run build`
- [ ] Run: `vercel --prod`
- [ ] Verify deployment at Vercel URL

### **Environment Variables** ⚠️

- [ ] Set `NEXT_PUBLIC_API_URL` = `https://your-backend.onrender.com`
- [ ] Set `NEXT_PUBLIC_AI_SERVICE_URL` = `https://your-ai-service.onrender.com`
- [ ] Verify in Vercel Dashboard → Settings → Environment Variables

---

## 🔧 **2. Backend → Render**

### **Service Creation** ⚠️

- [ ] Go to [Render Dashboard](https://dashboard.render.com)
- [ ] Create new Web Service
- [ ] Connect GitHub repository
- [ ] Use `render.yaml` from `backend/` directory

### **Environment Variables** ⚠️

- [ ] `ENVIRONMENT` = `production`
- [ ] `GIN_MODE` = `release`
- [ ] `PORT` = `10000`
- [ ] `JWT_SECRET` = `[generate: openssl rand -base64 64]`
- [ ] `CORS_ORIGIN` = `https://your-frontend.vercel.app`
- [ ] `MONGODB_URI` = `mongodb://your-mongo-host:27017/sms_app`
- [ ] `PLIVO_AUTH_ID` = `your_production_plivo_auth_id`
- [ ] `PLIVO_AUTH_TOKEN` = `your_production_plivo_auth_token`
- [ ] `PLIVO_FROM_NUMBER` = `your_production_plivo_phone_number`

### **Health Check** ✅

- [x] `/health` endpoint returns: `{"service":"sms-backend","status":"ok"}`

---

## 🤖 **3. AI Service → Render**

### **Service Creation** ⚠️

- [ ] Go to [Render Dashboard](https://dashboard.render.com)
- [ ] Create new Web Service
- [ ] Connect GitHub repository
- [ ] Use `render.yaml` from `ai-service/` directory

### **Environment Variables** ⚠️

- [ ] `ENVIRONMENT` = `production`
- [ ] `PORT` = `10000`
- [ ] `OPENAI_API_KEY` = `your_production_openai_api_key`
- [ ] `MODEL_NAME` = `gpt-4`
- [ ] `CORS_ORIGINS` = `https://your-frontend.vercel.app,https://your-backend.onrender.com`

### **Health Check** ✅

- [x] `/health` endpoint returns: `{"status":"healthy","service":"sms-ai-service","version":"1.0.0"}`

---

## 🗄️ **4. MongoDB Configuration**

### **Option A: Local MongoDB (Exposed)** ⚠️

- [ ] MongoDB running on port 27017
- [ ] Network access configured: `bindIp: 0.0.0.0`
- [ ] Firewall allows external connections
- [ ] Connection string: `mongodb://your-local-ip:27017/sms_app`

### **Option B: MongoDB Atlas (Recommended)** ⚠️

- [ ] Create cluster at [MongoDB Atlas](https://mongodb.com/atlas)
- [ ] Create database user with read/write permissions
- [ ] Get connection string: `mongodb+srv://user:pass@cluster.mongodb.net/sms_app`
- [ ] Add IP whitelist for Render services

---

## 🔐 **5. Security & CORS**

### **CORS Configuration** ✅

- [x] Backend CORS updated for production origins
- [x] AI service CORS updated for production origins
- [x] Frontend CORS origins configured in Vercel

### **Environment Security** ⚠️

- [ ] All sensitive keys in production environment variables
- [ ] No `.env` files committed to Git
- [ ] JWT secret is cryptographically secure
- [ ] HTTPS enforced everywhere

---

## 🧪 **6. Production Testing**

### **Health Checks** ⚠️

- [ ] Test: `curl https://your-backend.onrender.com/health`
- [ ] Test: `curl https://your-ai-service.onrender.com/health`
- [ ] Test: `curl https://your-frontend.vercel.app`

### **API Endpoints** ⚠️

- [ ] Test OTP: `POST /api/sms/send-otp`
- [ ] Test AI: `POST /chat`
- [ ] Test logs: `GET /api/logs`

### **Integration Testing** ⚠️

- [ ] Frontend → Backend communication
- [ ] Frontend → AI Service communication
- [ ] Backend → MongoDB connection
- [ ] Complete user flow: OTP → Verify → Dashboard → AI Chat

---

## 📊 **7. Monitoring & Logs**

### **Render Monitoring** ⚠️

- [ ] Enable log forwarding (optional)
- [ ] Set up health check alerts
- [ ] Monitor service performance

### **Vercel Analytics** ⚠️

- [ ] Enable Vercel Analytics
- [ ] Monitor frontend performance
- [ ] Track user experience metrics

---

## 🚨 **8. Troubleshooting**

### **Common Issues** ⚠️

- [ ] CORS errors: Check allowed origins in both services
- [ ] Database connection: Verify MongoDB URI and network access
- [ ] Build failures: Test builds locally first
- [ ] Port issues: Render uses `$PORT` environment variable

### **Debug Commands** ✅

```bash
# Service status
curl -v https://your-service.onrender.com/health

# Database connection
mongosh "your-mongodb-uri"

# Environment variables
echo $VARIABLE_NAME
```

---

## 🎯 **Final Verification**

### **All Systems Operational** ⚠️

- [ ] Frontend accessible at Vercel URL
- [ ] Backend API responding at Render URL
- [ ] AI service responding at Render URL
- [ ] MongoDB connection established
- [ ] Complete user flow working
- [ ] No CORS errors in browser console
- [ ] All environment variables set correctly

---

## 📚 **Resources**

- **Complete Guide**: [DEPLOYMENT.md](./DEPLOYMENT.md)
- **Quick Summary**: [DEPLOYMENT_SUMMARY.md](./DEPLOYMENT_SUMMARY.md)
- **Auto-Deploy Script**: [deploy.sh](./deploy.sh)
- **Docker Setup**: [docker-compose.yml](./docker-compose.yml)

---

**Status: Ready for Production Deployment! 🚀**

**Next Steps:**

1. Run `./deploy.sh` for automated deployment preparation
2. Follow Render deployment steps
3. Test all endpoints in production
4. Monitor service health and performance
