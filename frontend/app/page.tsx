'use client';

import { useEffect, useRef, useState } from 'react';
import Hls from 'hls.js';
import VideoPlayer from '@/components/VideoPlayer';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export default function Home() {
  const [videos, setVideos] = useState<string[]>([]);
  const [liveStreams, setLiveStreams] = useState<string[]>([]);
  const [selectedVideo, setSelectedVideo] = useState<string | null>(null);
  const [selectedLive, setSelectedLive] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [loadingLive, setLoadingLive] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadMessage, setUploadMessage] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    fetchVideos();
    fetchLiveStreams();
  }, []);

  const fetchVideos = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/videos`);
      if (!response.ok) {
        throw new Error('Failed to fetch videos');
      }
      const data = await response.json();
      setVideos(data.videos || []);
      if (data.videos && data.videos.length > 0) {
        setSelectedVideo(data.videos[0]);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load videos');
      console.error('Error fetching videos:', err);
    } finally {
      setLoading(false);
    }
  };

  const fetchLiveStreams = async () => {
    try {
      setLoadingLive(true);
      const response = await fetch(`${API_BASE_URL}/api/live`);
      if (!response.ok) {
        throw new Error('Failed to fetch live streams');
      }
      const data = await response.json();
      setLiveStreams(data.live || []);
    } catch (err) {
      console.error('Error fetching live streams:', err);
      // keep non-blocking; no user-facing error unless needed
    } finally {
      setLoadingLive(false);
    }
  };

  const handleUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    setError(null);
    setUploadMessage(null);
    setUploading(true);

    const name = file.name.replace(/\.[^/.]+$/, '');
    const formData = new FormData();
    formData.append('file', file);
    formData.append('name', name);

    try {
      const response = await fetch(`${API_BASE_URL}/api/videos/upload`, {
        method: 'POST',
        body: formData,
      });

      if (!response.ok) {
        const msg = await response.text();
        throw new Error(msg || 'Upload failed');
      }

      setUploadMessage(`Uploaded "${name}" successfully. Processing completed.`);
      await fetchVideos();
      setSelectedVideo(name);
      setSelectedLive(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
    }
  };

  return (
    <div className="container">
      <h1>Video Streaming Player</h1>

      {error && <div className="error">Error: {error}</div>}
      {uploadMessage && (
        <div className="notice success">
          {uploadMessage}
        </div>
      )}

      <div className="layout">
        <aside className="sidebar">
          <div className="upload-card">
            <h2>Upload a video</h2>
            <p>MP4/MOV/etc. We’ll convert to HLS automatically.</p>
            <input
              ref={fileInputRef}
              type="file"
              accept="video/*"
              onChange={handleUpload}
              disabled={uploading}
            />
            {uploading && <div className="loading">Uploading & converting...</div>}
          </div>

          <div className="video-list">
            <h2>Available Videos</h2>
            {loading ? (
              <div className="loading">Loading videos...</div>
            ) : videos.length === 0 ? (
              <p>No videos yet. Upload to get started.</p>
            ) : (
              videos.map((video) => (
                <div
                  key={video}
                  className={`video-item ${selectedVideo === video ? 'active' : ''}`}
                  onClick={() => {
                    setSelectedVideo(video);
                    setSelectedLive(null);
                  }}
                >
                  {video}
                </div>
              ))
            )}
          </div>

          <div className="video-list" style={{ marginTop: '1.5rem' }}>
            <h2>Live Streams</h2>
            {loadingLive ? (
              <div className="loading">Loading live streams...</div>
            ) : liveStreams.length === 0 ? (
              <p>No live streams active.</p>
            ) : (
              liveStreams.map((live) => (
                <div
                  key={live}
                  className={`video-item ${selectedLive === live ? 'active' : ''}`}
                  onClick={() => {
                    setSelectedLive(live);
                    setSelectedVideo(null);
                  }}
                >
                  {live}
                </div>
              ))
            )}
          </div>
        </aside>

        <main className="main-panel">
          {selectedVideo && (
            <VideoPlayer videoName={selectedVideo} />
          )}
          {selectedLive && !selectedVideo && (
            <VideoPlayer videoName={selectedLive} isLive />
          )}
          {!selectedVideo && !selectedLive && (
            <div className="empty-state">
              {loading ? 'Loading...' : 'Select a video or live stream to start playing.'}
            </div>
          )}
        </main>
      </div>
    </div>
  );
}

