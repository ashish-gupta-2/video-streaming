# Video Streaming Application

A full-stack video streaming application with Go backend and Next.js frontend, supporting HLS (HTTP Live Streaming) for recorded videos.

## Project Structure

```
video-streaming/
├── video-backend/          # Go backend server
│   ├── main.go
│   ├── handlers/
│   │   └── video.go
│   ├── videos/            # Video files directory
│   └── go.mod
└── frontend/              # Next.js frontend
    ├── app/
    ├── components/
    └── package.json
```

## Features

- ✅ HLS video streaming (Recorded videos)
- ✅ Video playlist management
- ✅ Responsive UI
- ✅ CORS support
- 🔄 Live streaming (Coming next)

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- FFmpeg (for video conversion)

### Backend Setup

1. Navigate to the backend directory:
```bash
cd video-backend
```

2. Install Go dependencies:
```bash
go mod download
```

3. Convert a video to HLS format (example):
```bash
mkdir -p videos/sample
ffmpeg -i your-video.mp4 \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/sample/segment%03d.ts" \
  videos/sample/index.m3u8
```

4. Start the backend server:
```bash
go run main.go
```

The backend will run on `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

The frontend will run on `http://localhost:3000`

## Usage

1. Start the backend server (port 8080)
2. Start the frontend server (port 3000)
3. Open `http://localhost:3000` in your browser
4. Select a video from the list to start streaming

## Converting Videos to HLS

To add new videos:

1. Create a directory for your video:
```bash
mkdir -p video-backend/videos/my-video
```

2. Convert the video using FFmpeg:
```bash
ffmpeg -i input.mp4 \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "video-backend/videos/my-video/segment%03d.ts" \
  video-backend/videos/my-video/index.m3u8
```

3. The video will automatically appear in the frontend video list

## Next Steps: Live Streaming

For live streaming, you'll need to:
1. Set up an RTMP server or use FFmpeg to stream to HLS
2. Implement real-time segment generation
3. Update the backend to handle live streams
4. Add live stream indicators in the frontend

## API Endpoints

- `GET /api/videos` - List all available videos
- `GET /api/videos/{videoName}/stream` - Get HLS playlist
- `GET /api/videos/{videoName}/segment/{segment}` - Get video segment
- `GET /health` - Health check

## License

MIT

