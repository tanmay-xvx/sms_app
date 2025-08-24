# 🚨 AI Service Deployment Troubleshooting

## ❌ **Common Build Errors & Solutions**

### **1. Metadata Generation Failed**

**Error:**

```
error: metadata-generation-failed
× Encountered error while generating package metadata.
```

**Solution:**

- Use `requirements-deploy.txt` instead of `requirements.txt`
- This file contains only essential dependencies
- Removes problematic packages like `redis`, `celery`, `python-multipart`

### **2. Python Version Issues**

**Error:**

```
Python version not supported
```

**Solution:**

- `runtime.txt` specifies Python 3.11.7
- Render supports Python 3.7-3.11
- Avoid Python 3.12+ for better compatibility

### **3. Pip Version Issues**

**Error:**

```
pip version outdated
```

**Solution:**

- Build command includes: `pip install --upgrade pip`
- This ensures latest pip version before installing packages

---

## 🔧 **Deployment Configuration**

### **Use Minimal Requirements**

```yaml
# render.yaml
buildCommand: |
  pip install --upgrade pip
  pip install -r requirements-deploy.txt
```

### **Requirements Files**

- **`requirements.txt`**: Full development dependencies
- **`requirements-deploy.txt`**: Minimal production dependencies
- **`runtime.txt`**: Python version specification

---

## 🧪 **Pre-Deployment Testing**

### **Test Locally First**

```bash
cd ai-service
pip install -r requirements-deploy.txt
python -c "import fastapi, uvicorn; print('All imports successful')"
```

### **Test Service Startup**

```bash
uvicorn main:app --host 0.0.0.0 --port 8000
```

---

## 📋 **Deployment Checklist**

- [ ] Use `requirements-deploy.txt` in render.yaml
- [ ] Verify `runtime.txt` specifies Python 3.11.7
- [ ] Test minimal requirements locally
- [ ] Ensure no problematic imports in main.py
- [ ] Set all required environment variables

---

## 🚀 **Quick Fix Commands**

### **If Build Still Fails**

```bash
# 1. Update render.yaml to use minimal requirements
buildCommand: |
  pip install --upgrade pip
  pip install fastapi uvicorn python-dotenv pydantic httpx openai

# 2. Or use specific versions
buildCommand: |
  pip install --upgrade pip
  pip install fastapi==0.104.1 uvicorn==0.24.0 python-dotenv==1.0.0
```

### **Environment Variables Required**

```bash
ENVIRONMENT=production
PORT=10000
OPENAI_API_KEY=your_key_here
MODEL_NAME=gpt-4
CORS_ORIGINS=https://your-frontend.vercel.app,https://your-backend.onrender.com
```

---

## 📚 **Alternative Solutions**

### **Option 1: Use requirements-deploy.txt** ✅

- Minimal dependencies
- Faster builds
- Better compatibility

### **Option 2: Inline Dependencies**

```yaml
buildCommand: |
  pip install --upgrade pip
  pip install fastapi uvicorn python-dotenv pydantic httpx openai
```

### **Option 3: Docker Deployment**

- Use Dockerfile instead of requirements
- More control over environment
- Consistent builds

---

**Status: Ready for deployment with minimal requirements! 🚀**
