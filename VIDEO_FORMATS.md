# Supported Video Formats

## Streaming Format

The system streams videos in **HLS (HTTP Live Streaming)** format:
- **Playlist file**: `.m3u8` (M3U8 playlist)
- **Video segments**: `.ts` (MPEG Transport Stream) files

## Input Video Formats (Source Videos)

You can convert **any video format that FFmpeg supports** to HLS. Common input formats include:

### Video Containers
- ✅ **MP4** (`.mp4`) - Most common
- ✅ **MOV** (`.mov`) - QuickTime
- ✅ **AVI** (`.avi`)
- ✅ **MKV** (`.mkv`) - Matroska
- ✅ **WebM** (`.webm`)
- ✅ **FLV** (`.flv`) - Flash Video
- ✅ **WMV** (`.wmv`) - Windows Media
- ✅ **3GP** (`.3gp`) - Mobile format
- ✅ **MPEG** (`.mpeg`, `.mpg`)
- ✅ **TS** (`.ts`) - Transport Stream

### Video Codecs (that FFmpeg can decode)
- H.264/AVC
- H.265/HEVC
- VP8, VP9
- AV1
- MPEG-2, MPEG-4
- And many more...

### Audio Codecs (that FFmpeg can decode)
- AAC
- MP3
- Opus
- Vorbis
- AC3
- And many more...

## Output Format (What Gets Streamed)

After conversion, all videos are streamed as **HLS** with:

### Video Codec
- **H.264** (recommended for best compatibility)
- **H.265/HEVC** (better compression, newer browsers)
- **VP9** (alternative, less common)

### Audio Codec
- **AAC** (recommended for best compatibility)
- **MP3** (alternative)

### Container
- **MPEG Transport Stream** (`.ts` files)

## Browser Compatibility

HLS streaming works on:
- ✅ **Safari** (native support)
- ✅ **Chrome** (via HLS.js library - already included)
- ✅ **Firefox** (via HLS.js library - already included)
- ✅ **Edge** (via HLS.js library - already included)
- ✅ **Mobile browsers** (iOS Safari, Chrome Mobile, etc.)

## Conversion Examples

### From MP4 to HLS
```bash
ffmpeg -i video.mp4 \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/segment%03d.ts" \
  videos/my-video/index.m3u8
```

### From MOV to HLS
```bash
ffmpeg -i video.mov \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/segment%03d.ts" \
  videos/my-video/index.m3u8
```

### From MKV to HLS
```bash
ffmpeg -i video.mkv \
  -c:v libx264 -c:a aac \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/segment%03d.ts" \
  videos/my-video/index.m3u8
```

### High Quality Conversion
```bash
ffmpeg -i input.mp4 \
  -c:v libx264 \
  -preset slow \
  -crf 22 \
  -c:a aac \
  -b:a 192k \
  -hls_time 10 \
  -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/segment%03d.ts" \
  videos/my-video/index.m3u8
```

### Multiple Quality Levels (Adaptive Bitrate)
```bash
# Generate multiple resolutions for adaptive streaming
ffmpeg -i input.mp4 \
  -c:v libx264 -c:a aac \
  -b:v 1000k -b:a 128k \
  -hls_time 10 -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/360p/segment%03d.ts" \
  videos/my-video/360p/index.m3u8

ffmpeg -i input.mp4 \
  -c:v libx264 -c:a aac \
  -b:v 2500k -b:a 192k \
  -hls_time 10 -hls_list_size 0 \
  -hls_segment_filename "videos/my-video/720p/segment%03d.ts" \
  videos/my-video/720p/index.m3u8
```

## Summary

**Input**: Any video format FFmpeg supports (MP4, MOV, AVI, MKV, WebM, etc.)  
**Output**: HLS format (.m3u8 + .ts segments)  
**Video Codec**: H.264 (recommended) or H.265  
**Audio Codec**: AAC (recommended) or MP3  
**Browser Support**: All modern browsers (via HLS.js or native)

The system is format-agnostic for input - as long as FFmpeg can read it, you can convert it to HLS and stream it!

