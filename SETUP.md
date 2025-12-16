# Setup Instructions

## Quick Start

### 1. Backend Setup

```bash
cd video-backend
go mod tidy
go run main.go
```

The backend will start on `http://localhost:8080`

### 2. Frontend Setup

In a new terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend will start on `http://localhost:3000`

### 3. Add a Sample Video

You need to convert a video to HLS format. Install FFmpeg if you haven't:

**macOS:**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt-get install ffmpeg
```

**Windows:**
Download from https://ffmpeg.org/download.html

Then convert a video:

```bash
# Create directory
mkdir -p video-backend/videos/sample

# Convert video (replace input.mp4 with your video file)
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "video-backend/videos/sample/segment%03d.ts" \
  video-backend/videos/sample/index.m3u8
```

### 4. Test the Application

1. Open `http://localhost:3000` in your browser
2. You should see your video in the list
3. Click on it to start streaming

## Troubleshooting

### Backend won't start
- Make sure Go is installed: `go version`
- Run `go mod tidy` to download dependencies
- Check if port 8080 is available

### Frontend won't start
- Make sure Node.js is installed: `node --version`
- Delete `node_modules` and run `npm install` again
- Check if port 3000 is available

### Video won't play
- Make sure the video is converted to HLS format correctly
- Check browser console for errors
- Verify the backend is running and accessible
- Try opening the playlist URL directly: `http://localhost:8080/api/videos/sample/stream`

### CORS errors
- The backend includes CORS headers, but if you see errors, check that the backend is running
- Make sure the API URL in the frontend matches the backend URL

## Next: Live Streaming

Once you have recorded video streaming working, we can implement live streaming by:
1. Setting up RTMP ingestion
2. Converting RTMP to HLS in real-time
3. Updating the backend to handle live streams
4. Adding live indicators in the frontend

