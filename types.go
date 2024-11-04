// types.go
package main

type Response struct {
	UUID     string  `json:"uuid"`
	Duration float64 `json:"duration"`
}

type TranscribedText struct {
	Done          bool      `json:"done"`
	Transcription string    `json:"transcription"`
	Duration      float32   `json:"duration"`
	Timestamps    []float64 `json:"timestamps"`
}

// type VideoRequest struct {
// 	VideoPath string `json:"videoPath"`
// }

// type VideoResponse struct {
// 	ProcessedVideoUrl string `json:"processedVideoUrl"`
// }
