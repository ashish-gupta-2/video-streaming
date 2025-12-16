# Deploy backend (Go) on Render (free)


1. Go to Render [dashboard](https://render.com/docs/free) → New + → Web Service.
2. Connect your GitHub account and select the video-streaming repo.
3. Settings:
    - Root Directory: video-backend
    - Environment: Go
    - Build Command: go build -o server .
    - Start Command: ./server
    - Region: closest to you
    - Instance Type: Free
4. Click Create Web Service and wait for deploy.
5. After deploy, note the URL, e.g. https://video-backend-xxxx.onrender.com.

That becomes API base: https://video-backend-xxxx.onrender.com.

## Update frontend .env.local:
```
cd frontend
echo "NEXT_PUBLIC_API_URL=https://video-backend-xxxx.onrender.com" > .env.local
git add frontend/.env.local
git commit -m "Point frontend to Render backend"
git push
cd frontendecho "NEXT_PUBLIC_API_URL=https://video-backend-xxxx.onrender.com" > .env.localgit add frontend/.env.localgit commit -m "Point frontend to Render backend"git push
```

# Deploy frontend on Vercel (free)

1. Go to Vercel → New Project → import the same repo.
2. Framework: auto-detected Next.js.
3. Root Directory: frontend
4. Add environment variable:
    - NEXT_PUBLIC_API_URL = https://video-backend-xxxx.onrender.com
5. Keep defaults, click Deploy.

After it builds, you’ll get a URL like https://video-streaming-xxxx.vercel.app.

# Storage & FFmpeg on cloud
- The current backend expects videos in ./videos, and live streams use FFmpeg running alongside the backend.
- On free cloud:
    - VOD uploads will work if:
        - Render service has writable disk (ephemeral; good for demos, not for long-term storage).
- Live streaming:
    - You’d run FFmpeg from your local machine, but point the -hls_segment_filename path to a cloud volume or S3-like store (needs more setup), or run FFmpeg in another Render service that shares storage (paid / advanced).

For an MVP demo:
- VOD upload & playback will work fine on Render.
- Live streaming from your camera is best kept local for now (as you have it) unless you want to add cloud storage + more infra.

