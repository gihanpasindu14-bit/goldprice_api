# Railway Deployment Guide

Complete guide to deploy the Gold Price API to Railway with a different Git account.

## Step 1: Switch Git Account Locally

### Option A: Configure for This Repository Only

```powershell
cd c:\Users\ASUS\goldprice_api

# Set the new git user for this repo only
git config user.name "Your Name"
git config user.email "your-new-email@example.com"

# Verify the change
git config user.name
git config user.email

# Check that it's local (not global)
git config --local user.name
```

### Option B: Configure Globally for All Repositories

```powershell
# Set globally
git config --global user.name "Your Name"
git config --global user.email "your-new-email@example.com"

# Verify
git config --global user.name
git config --global user.email
```

### Option C: Use SSH Key for Different Account

If using SSH authentication:

```powershell
# Generate new SSH key for the new account
ssh-keygen -t ed25519 -C "your-new-email@example.com" -f "C:\Users\ASUS\.ssh\id_railway"

# Add to ssh-agent
ssh-add C:\Users\ASUS\.ssh\id_railway

# Add public key to your GitHub account at github.com/settings/keys
```

## Step 2: Update Remote URL (If Needed)

If moving to a different GitHub account/repository:

```powershell
cd c:\Users\ASUS\goldprice_api

# Check current remote
git remote -v

# Change remote URL (replace USERNAME and REPO)
git remote set-url origin https://github.com/YOUR-NEW-USERNAME/goldprice_api.git

# Or for SSH
git remote set-url origin git@github.com:YOUR-NEW-USERNAME/goldprice_api.git

# Verify
git remote -v
```

## Step 3: Prepare Repository

```powershell
cd c:\Users\ASUS\goldprice_api

# Create commit with current changes
git add .
git commit -m "Prepare for Railway deployment with new git account"

# Or just stage updated files
git add Dockerfile .env.example .gitignore .railwayignore
git commit -m "Update deployment configuration for Railway"

# View git log to confirm
git log --oneline -5
```

## Step 4: Set Up Railway

### A. Create Railway Account & Project

1. Go to [railway.app](https://railway.app)
2. Sign up / Log in
3. Create a new project
4. Choose "Deploy from GitHub repo"
5. Authorize Railway to access your GitHub account
6. Select your repository (if it's not listed, you may need to push to new repo first)

### B. Push to New Repository

If you created a new repository:

```powershell
cd c:\Users\ASUS\goldprice_api

# Add new remote
git remote add railway https://github.com/YOUR-USERNAME/goldprice_api.git

# Push to new repository
git push -u railway main

# Or replace origin
git remote set-url origin https://github.com/YOUR-USERNAME/goldprice_api.git
git push -u origin main
```

## Step 5: Configure Environment Variables in Railway

1. In Railway dashboard, go to your project
2. Click on the service/deployment
3. Go to "Variables" tab
4. Add the following variables:

### Required Variables:

```env
PORT=8080
FIREBASE_CREDENTIALS_JSON={"type":"service_account","project_id":"..."}
```

### Get Your Firebase Credentials:

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Go to Project Settings → Service Accounts
4. Click "Generate New Private Key"
5. Copy the entire JSON content
6. Paste it as `FIREBASE_CREDENTIALS_JSON` variable in Railway

### Optional Variables:

```env
DEBUG=false
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://another-domain.com
```

## Step 6: Deploy to Railway

### Option A: Automatic Deployment (Recommended)

Railway will automatically deploy when you push to your connected GitHub repository:

```powershell
# Make changes
# Commit
git add .
git commit -m "Update feature"

# Push to GitHub
git push origin main

# Railway will automatically detect and deploy
```

### Option B: Manual Deployment via Railway CLI

```powershell
# Install Railway CLI
npm install -g @railway/cli

# Login to Railway
railway login

# Link to your project
railway link

# Deploy
railway up

# View logs
railway logs
```

## Step 7: Test Your Deployment

Once deployed, Railway will provide a public URL. Test it:

```powershell
# Test health check
curl https://your-app.railway.app/

# Test API
curl https://your-app.railway.app/api/metadata

# Upload CSV (with Postman or curl)
curl -X POST https://your-app.railway.app/api/upload \
  -F "file=@path/to/file.csv"
```

## Step 8: Connect Frontend to API

Update your frontend configuration in `goldprice_lk/js/main.js`:

```javascript
// Change from localhost to your Railway URL
const API_BASE_URL = 'https://your-app.railway.app';

// Or set via environment variable for flexibility
const API_BASE_URL = process.env.VUE_APP_API_URL || 'https://your-app.railway.app';
```

## Troubleshooting

### Issue: "Permission denied" on git push

**Solution:** 
- Check Git credentials: `git config user.name`
- If using HTTPS, you may need to use a [Personal Access Token](https://github.com/settings/tokens)
- For SSH, ensure your key is added to the new GitHub account

### Issue: Railway deployment fails with "firebase-credentials.json not found"

**Solution:**
- Don't copy `firebase-credentials.json` in Dockerfile (removed in updated version)
- Set `FIREBASE_CREDENTIALS_JSON` as environment variable in Railway instead

### Issue: CORS errors on frontend

**Solution:**
- Update `CORS_ALLOWED_ORIGINS` variable in Railway with your frontend domain
- Or update main.go to accept the environment variable (if not already implemented)

### Issue: Deployment stuck or slow

**Solution:**
- Check Railway logs: Go to project → Logs tab
- Wait 5-10 minutes for initial deployment
- Verify build logs for Go compilation errors

## Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `PORT` | Yes | 8080 | Port to run the API on |
| `FIREBASE_CREDENTIALS_JSON` | Yes | - | Firebase service account JSON (entire object as string) |
| `DEBUG` | No | false | Enable debug logging |
| `CORS_ALLOWED_ORIGINS` | No | * | Comma-separated list of allowed origins |

## Security Best Practices

1. ✅ Never commit `firebase-credentials.json`
2. ✅ Use environment variables for sensitive data
3. ✅ Set `DEBUG=false` in production
4. ✅ Restrict `CORS_ALLOWED_ORIGINS` to your domain
5. ✅ Use strong credentials in firebase-credentials.json
6. ✅ Rotate credentials periodically

## Next Steps

1. Deploy API to Railway
2. Update frontend to use Railway API URL
3. Test all endpoints
4. Monitor logs in Railway dashboard
5. Set up custom domain (optional)

## Useful Commands

```powershell
# View git status
git status

# View configured remotes
git remote -v

# View current branch
git branch

# Create new branch for deployment
git checkout -b deploy/railway

# Push specific branch
git push origin deploy/railway
```

## Support

- Railway Docs: [docs.railway.app](https://docs.railway.app)
- Firebase Admin SDK: [firebase.google.com/docs/admin](https://firebase.google.com/docs/admin)
- Go Gin Framework: [gin-gonic.com](https://gin-gonic.com)
