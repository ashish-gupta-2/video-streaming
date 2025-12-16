# Video Streaming Backend (Go)

Go backend server for streaming HLS (HTTP Live Streaming) videos.

## Prerequisites

- Go 1.21 or higher
- FFmpeg (for converting videos to HLS format)

## Installation

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run main.go
```

The server will start on port 8080 (or the port specified in the `PORT` environment variable).

## API Endpoints

- `GET /api/videos` - List all available videos
- `GET /api/videos/{videoName}/stream` - Get HLS playlist for a video
- `GET /api/videos/{videoName}/segment/{segment}` - Get a video segment
- `POST /api/videos/upload` - Upload a VOD file (auto-converts to HLS)
- `GET /api/live` - List live streams
- `GET /api/live/{streamName}/stream` - Get live HLS playlist
- `GET /api/live/{streamName}/segment/{segment}` - Get live segment
- `GET /health` - Health check endpoint

## Video Format

Videos must be in HLS format:
- Each video should be in its own directory under `./videos/`
- Each directory must contain:
  - `index.m3u8` - The HLS playlist file
  - `segment0.ts`, `segment1.ts`, etc. - Video segments

## Converting Videos to HLS

Use FFmpeg to convert your videos to HLS format:

```bash
# Basic conversion
ffmpeg -i input.mp4 -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls videos/sample/index.m3u8

# With better quality settings
ffmpeg -i input.mp4 \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/sample/segment%03d.ts" \
  videos/sample/index.m3u8
```

### Example:
```bash
# Create video directory
mkdir -p videos/my-video

# Convert video
ffmpeg -i my-video.mp4 \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/segment%03d.ts" \
  videos/my-video/index.m3u8
```

## Live Streaming (HLS)

We serve live streams from `videos/live/<streamName>/index.m3u8`. Use FFmpeg to continuously generate a live playlist and segments:

```bash
STREAM=event1
mkdir -p videos/live/$STREAM

ffmpeg -re -i input.mp4 \
  -c:v libx264 -c:a aac \
  -preset veryfast \
  -hls_time 2 \
  -hls_list_size 6 \
  -hls_flags delete_segments+append_list+program_date_time \
  -hls_segment_filename "videos/live/$STREAM/segment%03d.ts" \
  videos/live/$STREAM/index.m3u8
```

Then access:
- Playlist: `http://localhost:8080/api/live/$STREAM/stream`
- Segments are proxied automatically via the playlist

## Directory Structure

```
video-backend/
├── main.go
├── handlers/
│   └── video.go
├── videos/
│   └── sample/
│       ├── index.m3u8
│       ├── segment000.ts
│       ├── segment001.ts
│       └── ...
└── go.mod
```

## Environment Variables

- `PORT` - Server port (default: 8080)

