cd /Users/ashishgupta/Desktop/work/learn/video-streaming/video-backend

STREAM=cam1
mkdir -p videos/live/$STREAM

ffmpeg -re -f avfoundation -framerate 30 -pixel_format yuyv422 -video_size 1280x720 -i "0" \
  -f avfoundation -i ":0" \
  -c:v libx264 -preset veryfast -tune zerolatency -c:a aac \
  -g 30 -keyint_min 30 -sc_threshold 0 \
  -hls_time 1 -hls_list_size 4 \
  -hls_flags delete_segments+append_list+program_date_time \
  -hls_segment_filename "videos/live/$STREAM/segment%03d.ts" \
  videos/live/$STREAM/index.m3u8