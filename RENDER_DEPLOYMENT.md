# Deploy GoldPulse API on Render

## Prerequisites
- GitHub repository: https://github.com/gihanbandaraa/goldprice_api
- Firebase credentials JSON file
- Render account (free): https://render.com

## Step-by-Step Deployment

### 1. Sign Up / Login to Render
- Go to [render.com](https://render.com)
- Sign up with GitHub (easiest option)
- Authorize Render to access your repositories

### 2. Create New Web Service
1. Click **"New +"** → **"Web Service"**
2. Connect your GitHub account if not already connected
3. Select repository: **`gihanbandaraa/goldprice_api`**
4. Click **"Connect"**

### 3. Configure Service

**Basic Settings:**
- **Name:** `goldprice-api` (or your preferred name)
- **Region:** Choose closest to your users (e.g., Singapore for Asia)
- **Branch:** `main`
- **Root Directory:** Leave empty
- **Runtime:** `Go`
- **Build Command:** `go build -o main .`
- **Start Command:** `./main`

**Instance Type:**
- Select **"Free"** (automatically sleeps after 15min inactivity)
- Or **"Starter"** ($7/month - no sleep, better for production)

### 4. Add Environment Variables

Click **"Advanced"** → **"Add Environment Variable"**

Add these variables:

```
PORT=8000
ENVIRONMENT=production
GIN_MODE=release
ALLOWED_ORIGINS=*
```

**Important: Firebase Credentials**

Copy entire content from `firebase-credentials.json`:
```
FIREBASE_CREDENTIALS_JSON={"type":"service_account","project_id":"your-project",...}
```

### 5. Update Backend Code for Render

Your `services/firebase.go` already handles environment variables, so no code changes needed! ✅

### 6. Deploy

1. Click **"Create Web Service"**
2. Render will:
   - Clone your repository
   - Install Go dependencies
   - Build your application
   - Deploy it
3. Wait 2-3 minutes for deployment to complete

### 7. Get Your API URL

Once deployed, Render gives you a URL like:
```
https://goldprice-api.onrender.com
```

Copy this URL - you'll need it for the frontend!

### 8. Test Your API

Open in browser:
```
https://your-app-name.onrender.com/
https://your-app-name.onrender.com/api/prices/latest
```

Should see JSON responses ✅

## Update Frontend to Use Render API

### 1. Update API URL in Frontend

Edit `goldprice_lk/js/main.js`:

```javascript
const API_BASE_URL = 'https://goldprice-api.onrender.com';
```

### 2. Push Changes

```bash
cd goldprice_lk
git add .
git commit -m "Update API URL to Render"
git push origin main
```

Vercel will auto-redeploy ✅

### 3. Update CORS

Go back to Render dashboard:
1. Select your service
2. Go to **"Environment"** tab
3. Update `ALLOWED_ORIGINS` to your Vercel URL:
```
ALLOWED_ORIGINS=https://your-app.vercel.app
```
4. Click **"Save Changes"**
5. Service will auto-redeploy

## Important Notes

### Free Tier Limitations
- ⚠️ **Spins down after 15 minutes of inactivity**
- ⚠️ **Cold start takes 30-60 seconds** on first request
- ✅ 750 hours/month (enough for hobby projects)
- ✅ Auto-deploy on every git push

### Avoiding Cold Starts (Optional)

**Option 1: Upgrade to Starter Plan** ($7/month)
- No sleep/cold starts
- Better for production
- Predictable performance

**Option 2: Keep-Alive Service** (Free)
Use a cron job to ping your API every 10 minutes:
- UptimeRobot (free): https://uptimerobot.com
- Cron-job.org (free): https://cron-job.org

Add monitor for: `https://your-app.onrender.com/`

### Custom Domain (Optional)

1. Go to service **"Settings"** → **"Custom Domain"**
2. Add your domain (e.g., `api.goldpulse.lk`)
3. Update DNS records as instructed
4. SSL certificate auto-generated ✅

## Monitoring & Logs

### View Logs
1. Go to your service dashboard
2. Click **"Logs"** tab
3. See real-time logs from your Go application

### View Metrics
1. Click **"Metrics"** tab
2. See CPU, memory, response times

### Set Up Alerts (Optional)
1. Click **"Settings"** → **"Notifications"**
2. Add email for deploy failures/service issues

## Troubleshooting

### "Build Failed"
- Check logs for Go compilation errors
- Ensure `go.mod` and `go.sum` are committed
- Verify Go version compatibility

### "Service Unavailable"
- Check if service is sleeping (free tier)
- Wait 30-60 seconds for cold start
- Check logs for runtime errors

### "Firebase Connection Failed"
- Verify `FIREBASE_CREDENTIALS_JSON` is set correctly
- Check Firebase credentials are valid
- Ensure Firestore is enabled in Firebase Console

### "CORS Errors"
- Update `ALLOWED_ORIGINS` to include frontend URL
- Use exact URL (no trailing slash)
- Redeploy after changing environment variables

## Update Deployment

Render auto-deploys when you push to GitHub:

```bash
cd goldprice_api
# Make your changes
git add .
git commit -m "Your update message"
git push origin main
```

Render detects push and redeploys automatically ✅

## Manual Redeploy

In Render dashboard:
1. Go to your service
2. Click **"Manual Deploy"** → **"Deploy latest commit"**

## Costs

**Free Tier:**
- Free forever
- 750 hours/month
- Spins down after 15min inactivity
- **Perfect for development and testing**

**Starter Plan ($7/month):**
- Always on (no cold starts)
- Better performance
- **Recommended for production**

## Next Steps After Deployment

- [ ] Test all API endpoints
- [ ] Update frontend with Render URL
- [ ] Update CORS settings
- [ ] Set up uptime monitoring (optional)
- [ ] Consider upgrading to Starter plan for production
- [ ] Add custom domain (optional)

## Support

- Render Docs: https://render.com/docs
- Discord: https://render.com/discord
- Status: https://status.render.com

---

**Estimated deployment time: 5-10 minutes** ⚡
