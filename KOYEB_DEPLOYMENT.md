# Deploy GoldPulse API on Koyeb

## Why Koyeb?
- ✅ **No credit card required** for free tier
- ✅ No cold starts (always on)
- ✅ Auto-deploy from GitHub
- ✅ Free: 1 web service + 1 database
- ✅ Similar to Railway but truly free

## Prerequisites
- GitHub repository: https://github.com/gihanbandaraa/goldprice_api
- Firebase credentials JSON
- Koyeb account: https://koyeb.com (sign up with GitHub)

---

## Step-by-Step Deployment

### 1. Sign Up to Koyeb
1. Go to [koyeb.com](https://koyeb.com)
2. Click **"Sign up with GitHub"**
3. Authorize Koyeb to access repositories
4. **No credit card needed!** ✅

### 2. Create New App
1. Click **"Create App"**
2. Select **"GitHub"** as deployment method
3. Connect your GitHub account if not already
4. Select repository: **`gihanbandaraa/goldprice_api`**
5. Branch: **`main`**

### 3. Configure Deployment

**Builder:**
- Koyeb auto-detects Go - leave as **"Buildpack"**
- Build command: `go build -o main .`
- Run command: `./main`

**Instance:**
- Type: **Nano** (free tier)
- Region: Choose closest (e.g., Frankfurt, Singapore)

### 4. Set Environment Variables

Click **"Advanced"** → Add environment variables:

```env
PORT=8000
ENVIRONMENT=production
GIN_MODE=release
ALLOWED_ORIGINS=*
```

**Important: Firebase Credentials**

For Firebase, you have 2 options:

#### Option A: Environment Variable (Recommended)
1. Copy entire content from `firebase-credentials.json`
2. Add as single-line JSON:
```
FIREBASE_CREDENTIALS_JSON={"type":"service_account","project_id":"goldprice-api",...entire json...}
```

#### Option B: Upload as Secret
1. Click **"Secrets"** → **"Create Secret"**
2. Name: `firebase-credentials`
3. Type: **File**
4. Upload your `firebase-credentials.json`
5. Mount path: `/app/firebase-credentials.json`

### 5. Update Code for Secret File (if using Option B)

If using file secret, update `services/firebase.go`:

```go
func InitFirebase() (*firestore.Client, error) {
    ctx := context.Background()
    
    credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
    var opt option.ClientOption
    
    if credJSON != "" {
        // Option A: Environment variable
        opt = option.WithCredentialsJSON([]byte(credJSON))
    } else {
        // Option B: File secret
        credFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
        if credFile == "" {
            credFile = "firebase-credentials.json"
        }
        opt = option.WithCredentialsFile(credFile)
    }
    
    app, err := firebase.NewApp(ctx, nil, opt)
    // ... rest of code
}
```

But actually, your code already handles environment variable, so **use Option A** (no code changes needed) ✅

### 6. Health Check (Important!)

Koyeb needs a health check endpoint. Your `/` endpoint already works for this ✅

Set health check:
- Path: `/`
- Port: `8000`
- Protocol: `HTTP`

### 7. Deploy!

1. Review settings
2. Click **"Deploy"**
3. Wait 2-3 minutes for build & deployment
4. Watch logs in real-time

### 8. Get Your URL

Once deployed, Koyeb gives you a URL like:
```
https://goldprice-api-yourname.koyeb.app
```

### 9. Test Your API

```bash
# Health check
https://your-app.koyeb.app/

# Latest prices
https://your-app.koyeb.app/api/prices/latest
```

---

## Update Frontend

### 1. Update API URL

Edit `goldprice_lk/js/main.js`:

```javascript
const API_BASE_URL = 'https://goldprice-api-yourname.koyeb.app';
```

### 2. Push Changes

```bash
cd goldprice_lk
git add .
git commit -m "Update API URL to Koyeb"
git push origin main
```

### 3. Deploy Frontend on Vercel

- Go to [vercel.com](https://vercel.com)
- Import `goldprice_lk` repository
- Click "Deploy"
- Get Vercel URL

### 4. Update CORS

Go back to Koyeb:
1. Your app → **Settings** → **Environment Variables**
2. Edit `ALLOWED_ORIGINS`
3. Set to: `https://your-vercel-url.vercel.app`
4. Save → Auto-redeploys

---

## Important Notes

### Free Tier Benefits
- ✅ **No credit card required**
- ✅ **No cold starts** (always running)
- ✅ 1 web service free
- ✅ Auto-deploy on git push
- ✅ Built-in SSL/HTTPS
- ✅ DDoS protection

### Limitations
- 1 instance only (can't scale on free tier)
- Shared resources
- Limited regions
- 100 GB bandwidth/month

### Logs & Monitoring

**View Logs:**
1. Go to your app dashboard
2. Click **"Logs"** tab
3. See real-time application logs

**Metrics:**
1. Click **"Metrics"** tab
2. See CPU, memory, requests

### Auto-Deploy

Koyeb automatically redeploys when you push to GitHub:

```bash
# Make changes
git add .
git commit -m "Update"
git push origin main
```

Koyeb detects push and redeploys ✅

### Custom Domain (Optional)

1. Go to **Settings** → **Domains**
2. Click **"Add Custom Domain"**
3. Enter your domain (e.g., `api.goldpulse.lk`)
4. Update DNS records as instructed
5. SSL auto-configured ✅

---

## Troubleshooting

### Build Fails
- Check build logs in Koyeb dashboard
- Ensure `go.mod` and `go.sum` are committed
- Verify Go version compatibility

### App Crashes
- Check runtime logs
- Verify environment variables are set
- Test Firebase credentials locally first

### Can't Connect to Firestore
- Verify `FIREBASE_CREDENTIALS_JSON` is correctly formatted (single line)
- Check Firebase project is active
- Ensure Firestore is enabled

### CORS Errors
- Update `ALLOWED_ORIGINS` to exact Vercel URL
- No trailing slash in URL
- Redeploy after changing env vars

---

## Comparison: Koyeb vs Others

| Feature | Koyeb | Render | Railway |
|---------|-------|--------|---------|
| **Credit Card** | ❌ Not required | ✅ Required | ❌ Not required |
| **Cold Starts** | ❌ No | ✅ Yes (15min) | ❌ No |
| **Free Tier** | 1 service | 750 hrs | $5 credit |
| **Setup** | Easy | Easy | Easiest |
| **Best For** | No CC, always on | Hobby projects | Best UX |

**For your case: Koyeb is perfect!** ✅

---

## Cost Estimate

**Free Tier:**
- Free forever
- 1 web service
- 100 GB bandwidth
- No credit card needed
- **Perfect for this project** 🎉

**Paid Plans:**
- Start at $5/month
- More instances
- More resources

---

## Next Steps

1. ✅ Sign up on Koyeb (no CC)
2. ✅ Deploy backend from GitHub
3. ✅ Get Koyeb URL
4. ✅ Update frontend API URL
5. ✅ Deploy frontend on Vercel
6. ✅ Update CORS in Koyeb
7. ✅ Done!

**Estimated time: 5-10 minutes** ⚡

---

## Support

- Koyeb Docs: https://koyeb.com/docs
- Discord: https://koyeb.com/discord
- Status: https://status.koyeb.com

**Ready to deploy on Koyeb?** Start now at [koyeb.com](https://koyeb.com)!
