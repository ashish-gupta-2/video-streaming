'use client';

import { useEffect, useRef, useState } from 'react';
import Hls from 'hls.js';

interface VideoPlayerProps {
  videoName: string;
  isLive?: boolean;
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function VideoPlayer({ videoName, isLive = false }: VideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const hlsRef = useRef<Hls | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    const videoUrl = isLive
      ? `${API_BASE_URL}/api/live/${videoName}/stream`
      : `${API_BASE_URL}/api/videos/${videoName}/stream`;

    // Check if HLS is supported
    if (Hls.isSupported()) {
      const hls = new Hls({
        enableWorker: true,
        lowLatencyMode: false,
        backBufferLength: 90,
      });

      hlsRef.current = hls;

      hls.loadSource(videoUrl);
      hls.attachMedia(video);

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        setLoading(false);
        setError(null);
        video.play().catch((err) => {
          console.error('Error playing video:', err);
          setError('Failed to play video. User interaction may be required.');
        });
      });

      hls.on(Hls.Events.ERROR, (event, data) => {
        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              setError('Network error. Please check your connection.');
              hls.startLoad();
              break;
            case Hls.ErrorTypes.MEDIA_ERROR:
              setError('Media error. Trying to recover...');
              hls.recoverMediaError();
              break;
            default:
              setError('Fatal error. Cannot recover.');
              hls.destroy();
              break;
          }
        }
      });

      return () => {
        hls.destroy();
      };
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      // Native HLS support (Safari)
      video.src = videoUrl;
      video.addEventListener('loadedmetadata', () => {
        setLoading(false);
        setError(null);
      });
      video.addEventListener('error', () => {
        setError('Failed to load video');
        setLoading(false);
      });
    } else {
      setError('HLS is not supported in this browser');
      setLoading(false);
    }
  }, [videoName]);

  return (
    <div className="video-container">
      {loading && (
        <div style={{ padding: '2rem', textAlign: 'center', color: 'white' }}>
          Loading video...
        </div>
      )}
      {error && (
        <div style={{ padding: '1rem', background: '#fee', color: '#c33' }}>
          {error}
        </div>
      )}
      <video
        ref={videoRef}
        controls
        style={{ display: loading ? 'none' : 'block' }}
        playsInline
      />
    </div>
  );
}

