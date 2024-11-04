package main

import (
	//"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// var ReceivedVideoPath string

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or restrict to specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS) request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func handleVideoUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // Max 10 MB
	if err != nil {
		http.Error(w, "Could not parse form data", http.StatusBadRequest)
		return
	}

	// Get the uploaded video
	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Could not get uploaded video", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the overlay text
	text := r.FormValue("text")
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Save the video to a temporary file
	tempFile, err := os.CreateTemp("", handler.Filename)
	if err != nil {
		http.Error(w, "Could not save video", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // Clean up the temp file

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Could not save video", http.StatusInternalServerError)
		return
	}

	// Process the video with FFmpeg
	outputFile := tempFile.Name() + "_edited.mp4"
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(),
		"-vf", fmt.Sprintf("drawtext=text='%s':fontcolor=white:fontsize=24:x=(w-text_w)/2:y=h-(text_h+20)", text),
		"-codec:a", "copy", outputFile)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Failed to process video", http.StatusInternalServerError)
		return
	}
	defer os.Remove(outputFile)

	// Send the edited video back to the client
	w.Header().Set("Content-Type", "video/mp4")
	output, err := os.Open(outputFile)
	if err != nil {
		http.Error(w, "Could not open edited video", http.StatusInternalServerError)
		return
	}
	defer output.Close()

	_, err = io.Copy(w, output)
	if err != nil {
		http.Error(w, "Could not send video", http.StatusInternalServerError)
		return
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Limit your max memory to 10MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form data
	file, _, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file to save the uploaded video
	out, err := os.Create("uploaded_video.mp4") // Change path as necessary
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copy the uploaded file to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	err = addSubtitleToVideo("uploaded_video.mp4", "subtitles.srt", "output_with_subtitles.mp4")
	if err != nil {
		http.Error(w, "Failed to process video with subtitles", http.StatusInternalServerError)
		return
	}
	processedVideoUrl := fmt.Sprintf("http://localhost:8000/output_with_subtitles.mp4")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"processedVideoUrl": "%s"}`, processedVideoUrl)
}
