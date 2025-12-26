# Deploy GoldPulse API on Fly.io

## Why Fly.io?
- ✅ **100% FREE** - No credit card required for free tier
- ✅ **No cold starts** - Always running
- ✅ **Global deployment** - Edge locations worldwide
- ✅ **3 shared VMs free** - More than enough for this project
- ✅ **160 GB bandwidth/month** - Generous free tier

---

## Prerequisites
- GitHub repository: https://github.com/gihanbandaraa/goldprice_api
- Firebase credentials JSON
- Fly.io account (free, no CC)
- Fly CLI installed

---

## Step 1: Install Fly CLI

### Windows (PowerShell):
```powershell
iwr https://fly.io/install.ps1 -useb | iex
```

### macOS/Linux:
```bash
curl -L https://fly.io/install.sh | sh
```

After installation, restart your terminal.

---

## Step 2: Sign Up & Login

```bash
# Sign up (opens browser, use GitHub login)
fly auth signup

# Or login if you already have account
fly auth login
```

**No credit card required!** ✅

---

## Step 3: Navigate to Project

```bash
cd goldprice_api
```

---

## Step 4: Launch App on Fly.io

```bash
fly launch
```

This will:
1. Detect your Dockerfile
2. Ask configuration questions
3. Create `fly.toml` (already created for you)

**Answer the prompts:**
- App name: `goldprice-api` (or choose your own)
- Region: Choose closest (Singapore = `sin`, Frankfurt = `fra`)
- PostgreSQL database: **No** (we use Firebase)
- Redis database: **No**
- Deploy now: **No** (we need to set secrets first)

---

## Step 5: Set Environment Secrets

### Set Firebase Credentials:

```bash
# Copy your firebase-credentials.json content and set as secret
fly secrets set FIREBASE_CREDENTIALS_JSON='{"type":"service_account","project_id":"your-project-id",...}'
```

**Important:** Copy the ENTIRE firebase-credentials.json as a single line!

### Set CORS:

```bash
# Temporarily allow all origins
fly secrets set ALLOWED_ORIGINS='*'
```

You'll update this later with your Vercel URL.

---

## Step 6: Deploy!

```bash
fly deploy
```

This will:
1. Build Docker image
2. Push to Fly.io registry
3. Deploy to your VM
4. Start your application

**Wait 2-3 minutes** for deployment to complete.

---

## Step 7: Get Your URL

After deployment completes:

```bash
fly status
```

Your URL will be: `https://goldprice-api.fly.dev`

Or open in browser:
```bash
fly open
```

---

## Step 8: Test Your API

```bash
# Health check
curl https://goldprice-api.fly.dev/

# Latest prices
curl https://goldprice-api.fly.dev/api/prices/latest
```

Should return JSON responses ✅

---

## Step 9: View Logs

```bash
# Real-time logs
fly logs

# Or in dashboard
fly dashboard
```

---

## Update Frontend with Fly.io URL

### 1. Update API URL

Edit `goldprice_lk/js/main.js`:

```javascript
const API_BASE_URL = 'https://goldprice-api.fly.dev';
```

### 2. Commit and Push

```bash
cd goldprice_lk
git add .
git commit -m "Update API URL to Fly.io"
git push origin main
```

### 3. Deploy Frontend on Vercel

- Go to [vercel.com](https://vercel.com)
- Import `goldprice_lk` repository
- Click "Deploy"
- Get your Vercel URL

### 4. Update CORS

```bash
cd goldprice_api
fly secrets set ALLOWED_ORIGINS='https://your-app.vercel.app'
```

App will auto-restart with new CORS settings ✅

---

## Useful Fly.io Commands

```bash
# View app status
fly status

# View logs (real-time)
fly logs

# SSH into your VM
fly ssh console

# Open app in browser
fly open

# View dashboard
fly dashboard

# Scale app (free tier allows 1 machine)
fly scale count 1

# Restart app
fly apps restart goldprice-api

# Destroy app (be careful!)
fly apps destroy goldprice-api
```

---

## Update Your App

When you make code changes:

```bash
cd goldprice_api
git add .
git commit -m "Your changes"
git push origin main

# Deploy to Fly.io
fly deploy
```

---

## Free Tier Limits

- ✅ 3 shared-cpu VMs (1x-256mb) - **You're using 1**
- ✅ 160 GB outbound bandwidth/month
- ✅ 3 GB persistent volume storage
- ✅ Always running (no cold starts)
- ✅ Automatic SSL/HTTPS
- ✅ Global Anycast network

**Your app will stay FREE forever on this tier!** 🎉

---

## Monitoring

### Check App Health:
```bash
fly checks list
```

### View Metrics:
```bash
fly dashboard
```

Go to Metrics tab to see:
- CPU usage
- Memory usage
- Request count
- Response times

---

## Custom Domain (Optional)

### Add Custom Domain:

```bash
fly certs create api.goldpulse.lk
```

Then add DNS records (instructions shown after command).

SSL certificate auto-generated ✅

---

## Troubleshooting

### Build Fails
```bash
# Check build logs
fly logs

# Rebuild from scratch
fly deploy --build-only
```

### App Crashes
```bash
# View crash logs
fly logs --app goldprice-api

# Check if app is running
fly status

# SSH into VM to debug
fly ssh console
```

### Can't Connect to Firebase
```bash
# Verify secrets are set
fly secrets list

# Re-set Firebase credentials
fly secrets set FIREBASE_CREDENTIALS_JSON='...'
```

### CORS Issues
```bash
# Update CORS
fly secrets set ALLOWED_ORIGINS='https://your-vercel-url.vercel.app'

# Check if secret is set
fly secrets list
```

---

## Scaling (If Needed in Future)

```bash
# Add more VMs (costs $$ beyond free tier)
fly scale count 2

# Increase memory (costs $$)
fly scale memory 512

# Scale back to free tier
fly scale count 1
fly scale memory 256
```

---

## Regions

Fly.io has edge locations worldwide. Current setup: `sin` (Singapore)

Change region:
```bash
# List all regions
fly platform regions

# Add region (multi-region)
fly regions add fra
```

---

## Environment Variables

### View Current Secrets:
```bash
fly secrets list
```

### Set New Secret:
```bash
fly secrets set KEY=VALUE
```

### Remove Secret:
```bash
fly secrets unset KEY
```

App auto-restarts after changing secrets.

---

## Cost Estimate

**Free Tier (Your Current Setup):**
- 1 VM (shared-cpu, 256MB) = **$0/month**
- 160 GB bandwidth = **$0/month**
- SSL certificate = **$0/month**
- **Total: $0/month forever** 🎉

**If You Scale Up Later:**
- Additional VMs = ~$2/month each
- More memory = +$0.50/month per 128MB
- Extra bandwidth = $0.02/GB over 160GB

---

## Why Fly.io is Perfect for This Project

✅ **No Credit Card** - Sign up with GitHub, deploy for free
✅ **Always On** - No cold starts, instant responses
✅ **Global CDN** - Fast from anywhere
✅ **Simple Deployment** - One command: `fly deploy`
✅ **Docker-Based** - Full control over environment
✅ **Great Logs** - Easy debugging
✅ **Free SSL** - Automatic HTTPS
✅ **Perfect for Go** - Optimized for compiled apps

---

## Support

- Fly.io Docs: https://fly.io/docs
- Community: https://community.fly.io
- Discord: https://fly.io/discord
- Status: https://status.flyio.net

---

## Quick Reference

```bash
# Deploy
fly deploy

# Logs
fly logs

# Status
fly status

# Dashboard
fly dashboard

# Restart
fly apps restart

# Secrets
fly secrets set KEY=value
fly secrets list
```

---

## Next Steps After Deployment

- [ ] Test all API endpoints
- [ ] Update frontend with Fly.io URL
- [ ] Deploy frontend on Vercel
- [ ] Update CORS with Vercel URL
- [ ] Set up monitoring alerts (optional)
- [ ] Add custom domain (optional)

**Total deployment time: ~10 minutes** ⚡

**Your API will be live at:** `https://goldprice-api.fly.dev`
