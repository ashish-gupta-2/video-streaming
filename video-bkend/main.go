package main

import (
	"log"
	"net/http"
	"os"

	"video-backend/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Create videos directory if it doesn't exist
	if err := os.MkdirAll("./videos", 0755); err != nil {
		log.Fatal("Failed to create videos directory:", err)
	}
	// Create live videos directory
	if err := os.MkdirAll("./videos/live", 0755); err != nil {
		log.Fatal("Failed to create live videos directory:", err)
	}

	// Create router
	r := mux.NewRouter()

	// CORS middleware
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, req)
		})
	}

	// Apply CORS middleware
	r.Use(corsMiddleware)

	// Video streaming routes
	r.HandleFunc("/api/videos/{videoName}/stream", handlers.StreamVideo).Methods("GET")
	r.HandleFunc("/api/videos/{videoName}/segment/{segment}", handlers.ServeSegment).Methods("GET")
	r.HandleFunc("/api/videos", handlers.ListVideos).Methods("GET")
	r.HandleFunc("/api/videos/upload", handlers.UploadVideo).Methods("POST")
	// Live streaming routes
	r.HandleFunc("/api/live", handlers.ListLive).Methods("GET")
	r.HandleFunc("/api/live/{streamName}/stream", handlers.StreamLive).Methods("GET")
	r.HandleFunc("/api/live/{streamName}/segment/{segment}", handlers.ServeLiveSegment).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Video files should be placed in ./videos/ directory")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
