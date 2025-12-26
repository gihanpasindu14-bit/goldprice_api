# Deployment Guide

## Prerequisites

1. GitHub account
2. Railway account ([railway.app](https://railway.app))
3. Firebase project with Firestore enabled
4. Firebase service account credentials

## Railway Deployment (Backend API)

### Option 1: Deploy from GitHub (Recommended)

1. **Push code to GitHub**
```bash
cd goldprice_api
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/goldprice-api.git
git push -u origin main
```

2. **Deploy on Railway**
   - Go to [railway.app](https://railway.app)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository
   - Railway will auto-detect Go project

3. **Add Environment Variables**
   - In Railway dashboard, go to your project
   - Click "Variables" tab
   - Add these variables:
     ```
     PORT=8000
     ENVIRONMENT=production
     ALLOWED_ORIGINS=https://your-frontend-url.vercel.app
     ```

4. **Add Firebase Credentials**
   - Copy content from `firebase-credentials.json`
   - In Railway Variables, add:
     ```
     FIREBASE_CREDENTIALS_JSON=<paste entire JSON content>
     ```
   - Update `services/firebase.go` to read from env variable

5. **Get your Railway URL**
   - Railway will generate a URL like: `https://goldprice-api-production.up.railway.app`
   - Copy this URL for frontend configuration

### Option 2: Railway CLI

```bash
npm i -g @railway/cli
railway login
railway init
railway up
```

## Vercel Deployment (Frontend)

### Option 1: Deploy from GitHub

1. **Push frontend to GitHub**
```bash
cd goldprice_lk
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/goldprice-lk.git
git push -u origin main
```

2. **Deploy on Vercel**
   - Go to [vercel.com](https://vercel.com)
   - Click "Add New Project"
   - Import your GitHub repository
   - Vercel will auto-detect static site
   - Click "Deploy"

3. **Update API URL**
   - After Railway deployment, copy your API URL
   - Update `js/main.js`:
     ```javascript
     const API_BASE_URL = 'https://goldprice-api-production.up.railway.app';
     ```
   - Commit and push changes
   - Vercel will auto-redeploy

### Option 2: Vercel CLI

```bash
npm i -g vercel
cd goldprice_lk
vercel
```

## Update CORS on Backend

After deploying frontend, update Railway environment variable:

```
ALLOWED_ORIGINS=https://your-actual-vercel-url.vercel.app
```

Then update `main.go` CORS configuration to read from environment.

## Firebase Credentials on Railway

### Method 1: Environment Variable (Recommended)

Update `services/firebase.go`:

```go
func InitFirebase() (*firestore.Client, error) {
    ctx := context.Background()
    
    // Check if credentials are in environment variable
    credJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")
    var opt option.ClientOption
    
    if credJSON != "" {
        // Production: Use env variable
        opt = option.WithCredentialsJSON([]byte(credJSON))
    } else {
        // Development: Use file
        opt = option.WithCredentialsFile("firebase-credentials.json")
    }
    
    app, err := firebase.NewApp(ctx, nil, opt)
    if err != nil {
        return nil, err
    }
    
    client, err := app.Firestore(ctx)
    if err != nil {
        return nil, err
    }
    
    return client, nil
}
```

### Method 2: Railway Secret Files

1. In Railway dashboard → Settings → Secret Files
2. Add file: `firebase-credentials.json`
3. Paste your credentials JSON

## Post-Deployment Checklist

- [ ] Backend API is accessible at Railway URL
- [ ] Frontend is accessible at Vercel URL
- [ ] API health check works: `GET https://your-api.railway.app/`
- [ ] Latest prices endpoint works: `GET https://your-api.railway.app/api/prices/latest`
- [ ] Frontend can fetch data from API (check browser console)
- [ ] CORS is properly configured
- [ ] Firebase Firestore connection works
- [ ] All pages load correctly (Dashboard, Alerts, Calculators, etc.)

## Custom Domains (Optional)

### Vercel Custom Domain
1. Go to Project Settings → Domains
2. Add your custom domain (e.g., `goldpulse.lk`)
3. Update DNS records as instructed

### Railway Custom Domain
1. Go to Project Settings → Networking
2. Add custom domain (e.g., `api.goldpulse.lk`)
3. Update DNS records as instructed

## Monitoring

### Railway
- View logs: Railway Dashboard → Deployments → View Logs
- Monitor usage: Railway Dashboard → Metrics

### Vercel
- View analytics: Vercel Dashboard → Analytics
- Check deployment logs: Vercel Dashboard → Deployments

## Troubleshooting

### API Returns 500 Error
- Check Railway logs for errors
- Verify Firebase credentials are correctly set
- Ensure Firestore is enabled in Firebase Console

### CORS Errors in Browser
- Update `ALLOWED_ORIGINS` in Railway to include your Vercel URL
- Redeploy backend after changing CORS settings

### Frontend Shows No Data
- Check browser console for errors
- Verify API_BASE_URL in `js/main.js` is correct
- Test API endpoint directly in browser

### Firebase Connection Fails
- Verify Firebase credentials JSON is valid
- Check Firestore security rules
- Ensure Firebase project billing is enabled

## Costs

- **Vercel**: Free tier (100 GB bandwidth/month)
- **Railway**: Free $5 credit/month, then pay-as-you-go
- **Firebase**: Free tier (50K reads/day)

Estimated monthly cost: **$0-5** for small traffic

## Support

For deployment issues:
- Railway: [docs.railway.app](https://docs.railway.app)
- Vercel: [vercel.com/docs](https://vercel.com/docs)
- Firebase: [firebase.google.com/docs](https://firebase.google.com/docs)
