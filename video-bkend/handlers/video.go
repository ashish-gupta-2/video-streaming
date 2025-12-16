package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const videosDir = "./videos"
const liveDir = "./videos/live"

// StreamVideo serves the HLS playlist (m3u8 file)
func StreamVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoName := vars["videoName"]

	// Sanitize video name to prevent directory traversal
	videoName = filepath.Base(videoName)

	servePlaylist(videosDir, fmt.Sprintf("/api/videos/%s/segment/", videoName), videoName, w)
}

// ServeSegment serves individual video segments (.ts files)
func ServeSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoName := vars["videoName"]
	segment := vars["segment"]

	// Sanitize inputs
	videoName = filepath.Base(videoName)
	segment = filepath.Base(segment)

	// Ensure segment has .ts extension
	if !strings.HasSuffix(segment, ".ts") {
		segment = segment + ".ts"
	}

	segmentPath := filepath.Join(videosDir, videoName, segment)

	// Check if segment exists
	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		http.Error(w, "Segment not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers for video segments
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Open and serve the segment
	file, err := os.Open(segmentPath)
	if err != nil {
		log.Printf("Error opening segment: %v", err)
		http.Error(w, "Error serving segment", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Stream the file
	io.Copy(w, file)
}

// StreamLive serves a live HLS playlist under videos/live/<streamName>/index.m3u8
func StreamLive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamName := vars["streamName"]
	streamName = filepath.Base(streamName)

	servePlaylist(liveDir, fmt.Sprintf("/api/live/%s/segment/", streamName), streamName, w)
}

// ServeLiveSegment serves live .ts segments from videos/live/<streamName>
func ServeLiveSegment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	streamName := filepath.Base(vars["streamName"])
	segment := filepath.Base(vars["segment"])

	if !strings.HasSuffix(segment, ".ts") {
		segment = segment + ".ts"
	}

	segmentPath := filepath.Join(liveDir, streamName, segment)

	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		http.Error(w, "Segment not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// shorter cache for live
	w.Header().Set("Cache-Control", "public, max-age=15")

	file, err := os.Open(segmentPath)
	if err != nil {
		log.Printf("Error opening live segment: %v", err)
		http.Error(w, "Error serving segment", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	io.Copy(w, file)
}

// UploadVideo handles video uploads, converts them to HLS, and stores them under videosDir/<name>/
func UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Allow up to ~1GB uploads; adjust as needed
	if err := r.ParseMultipartForm(1 << 30); err != nil {
		http.Error(w, "failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "video file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine video name (folder name)
	userProvidedName := strings.TrimSpace(r.FormValue("name"))
	videoName := userProvidedName
	if videoName == "" {
		videoName = strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	}
	videoName = sanitizeName(videoName)
	if videoName == "" {
		videoName = fmt.Sprintf("upload-%d", time.Now().Unix())
	}

	targetDir := filepath.Join(videosDir, videoName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		http.Error(w, "failed to create video directory", http.StatusInternalServerError)
		return
	}

	// Save the uploaded file temporarily
	tmpInput := filepath.Join(targetDir, "source"+filepath.Ext(header.Filename))
	out, err := os.Create(tmpInput)
	if err != nil {
		http.Error(w, "failed to save upload", http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		http.Error(w, "failed to write upload", http.StatusInternalServerError)
		return
	}
	out.Close()

	// Ensure ffmpeg exists
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		http.Error(w, "ffmpeg not installed", http.StatusInternalServerError)
		return
	}

	// Build ffmpeg command to generate HLS
	playlistPath := filepath.Join(targetDir, "index.m3u8")
	segmentPattern := filepath.Join(targetDir, "segment%03d.ts")
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-i", tmpInput,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-hls_time", "6",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentPattern,
		playlistPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("ffmpeg error: %v, output: %s", err, string(output))
		http.Error(w, "failed to convert video to HLS", http.StatusInternalServerError)
		return
	}

	// Optional: keep the source; uncomment next line to remove
	// _ = os.Remove(tmpInput)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"video":"%s","playlist":"/api/videos/%s/stream"}`, videoName, videoName)))
}

// ListVideos returns a list of available videos
func ListVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var videos []string

	// Read videos directory
	entries, err := os.ReadDir(videosDir)
	if err != nil {
		log.Printf("Error reading videos directory: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to read videos directory"}`))
		return
	}

	// Check each entry for valid video directories (containing index.m3u8)
	for _, entry := range entries {
		if entry.IsDir() {
			playlistPath := filepath.Join(videosDir, entry.Name(), "index.m3u8")
			if _, err := os.Stat(playlistPath); err == nil {
				videos = append(videos, entry.Name())
			}
		}
	}

	// Return JSON response
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string][]string{
		"videos": videos,
	})
}

// ListLive returns a list of live streams (folders under videos/live with index.m3u8)
func ListLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var streams []string

	entries, err := os.ReadDir(liveDir)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string][]string{"live": streams})
			return
		}
		log.Printf("Error reading live directory: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read live directory"})
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			playlistPath := filepath.Join(liveDir, entry.Name(), "index.m3u8")
			if _, err := os.Stat(playlistPath); err == nil {
				streams = append(streams, entry.Name())
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string][]string{"live": streams})
}

// sanitizeName keeps alphanumerics, dash and underscore only
func sanitizeName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	clean := re.ReplaceAllString(name, "-")
	clean = strings.Trim(clean, "-")
	return clean
}

// servePlaylist reads index.m3u8, rewrites segment URLs to the API, and serves it.
func servePlaylist(rootDir, baseSegmentURL, name string, w http.ResponseWriter) {
	playlistPath := filepath.Join(rootDir, name, "index.m3u8")

	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// no-cache so clients pull the freshest playlist for live and vod
	w.Header().Set("Cache-Control", "no-cache")

	file, err := os.Open(playlistPath)
	if err != nil {
		log.Printf("Error opening playlist: %v", err)
		http.Error(w, "Error serving playlist", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	playlistContent, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading playlist: %v", err)
		http.Error(w, "Error reading playlist", http.StatusInternalServerError)
		return
	}

	lines := strings.Split(string(playlistContent), "\n")
	var modifiedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "segment") && strings.HasSuffix(trimmed, ".ts") {
			modifiedLines = append(modifiedLines, baseSegmentURL+trimmed)
		} else {
			modifiedLines = append(modifiedLines, line)
		}
	}

	modifiedContent := strings.Join(modifiedLines, "\n")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(modifiedContent))
}
