package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

var transcriptionData *TranscribedText // Global variable to store transcription data

func main() {
	c := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	http.Handle("/output_with_subtitles.mp4", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/api/message", messageHandler)

	// Start transcription process in a separate goroutine
	go startTranscription()

	fmt.Println("Server started at :8000")
	err := http.ListenAndServe(":8000", c(http.DefaultServeMux))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// startTranscription initiates the audio transcription process
func startTranscription() {
	filePath := "uploaded_video.mp4"

	// Wait until the file exists
	for {
		if _, err := os.Stat(filePath); err == nil {
			break // File exists
		} else if os.IsNotExist(err) {
			//fmt.Println("File does not exist yet, waiting...")
		} else {
			log.Fatalf("Error checking file existence: %v", err)
		}
		time.Sleep(1 * time.Second) // Check every second
	}

	uuid, err := pushAudio(filePath)
	if err != nil {
		log.Fatalf("Error pushing audio: %v", err)
	}

	for {
		data, err := getOutputText(uuid)
		if err != nil {
			log.Printf("Error fetching transcription: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if data.Done {
			transcriptionData = data
			jsonData, _ := json.MarshalIndent(data, "", "    ")
			fmt.Println("Full JSON response:", string(jsonData))
			createSubtitles() // Create subtitles after transcription is done
			break
		} else {
			fmt.Println("Transcription in progress. Retrying in 5 seconds...")
		}

		time.Sleep(5 * time.Second)
	}
}

// messageHandler handles requests to get transcription data
func messageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if transcriptionData != nil {
		json.NewEncoder(w).Encode(transcriptionData)
	} else {
		http.Error(w, "Transcription not available yet", http.StatusServiceUnavailable)
	}
}
