// step3.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func createSubtitles() {
	// Wait until transcriptionData is available
	for transcriptionData == nil {
		// fmt.Println("Хадмал текстийн мэдээлэл хүлээгдэж байна...")
		time.Sleep(1 * time.Second) // Wait for 1 second before checking again
	}

	transcription := transcriptionData.Transcription
	timestamps := transcriptionData.Timestamps
	words := strings.Split(transcription, " ")

	var filteredTimestamps []float64
	if len(timestamps) > 0 {
		filteredTimestamps = append(filteredTimestamps, timestamps[0]) // Эхний timestamp-г оруулна
		for i := 1; i < len(timestamps); i++ {
			// Өмнөх timestamp-тай ижил биш бол шинэ массивд нэмнэ
			if timestamps[i] != timestamps[i-1] {
				filteredTimestamps = append(filteredTimestamps, timestamps[i])
			}
		}
	}

	if len(words) > len(timestamps) {
		fmt.Println("Анхааруулга: Үгсийн тоо timestamp-үүдээс их байна.")
		words = words[:len(timestamps)]
	} else if len(words) < len(timestamps) {
		fmt.Println("Анхааруулга: Timestamp-үүдийн тоо үгсээс их байна.")
		timestamps = timestamps[:len(words)]
	}

	file, err := os.Create("subtitles.srt")
	if err != nil {
		fmt.Println("Алдаа: Хадмал текстийн файл үүсгэхдээ алдаа гарлаа:", err)
		return
	}
	defer file.Close()

	for i := 0; i < len(words) && i < len(filteredTimestamps); i++ {
		startTime := filteredTimestamps[i]
		endTime := startTime + 3.0

		if i < len(filteredTimestamps)-1 {
			endTime = filteredTimestamps[i+1]
		}

		start := formatTime(startTime)
		end := formatTime(endTime)

		if _, err := fmt.Fprintf(file, "%d\n%s --> %s\n%s\n\n", i+1, start, end, words[i]); err != nil {
			fmt.Println("Алдаа: Файлд бичих үед алдаа гарлаа:", err)
			return
		}
	}
	fmt.Println("Хадмал текстийн файл 'subtitles.srt' амжилттай үүсгэгдлээ!")

	inputVideo := "input.mp4"
	subtitleFile := "subtitles.srt"
	outputVideo := "output_with_subtitles.mp4"

	err = addSubtitleToVideo(inputVideo, subtitleFile, outputVideo)
	if err != nil {
		fmt.Println("Алдаа:", err)
	}

}

func addSubtitleToVideo(inputVideo, subtitleFile, outputVideo string) error {
	// Command to add subtitle using FFmpeg
	cmd := exec.Command("ffmpeg", "-i", inputVideo, "-vf", fmt.Sprintf("subtitles=%s", subtitleFile), outputVideo)

	// Run the FFmpeg command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add subtitle to video: %w", err)
	}

	fmt.Println("Subtitled video created successfully at", outputVideo)
	return nil
}

func formatTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) / 60) % 60
	secondsInt := int(seconds) % 60
	milliseconds := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secondsInt, milliseconds)
}

// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"strings"
// )

// func createSubtitles(inputVideo string) {
// 	if transcriptionData == nil {
// 		fmt.Println("Алдаа: Хадмал текстийн мэдээлэл байхгүй.")
// 		return
// 	}

// 	transcription := transcriptionData.Transcription
// 	timestamps := transcriptionData.Timestamps
// 	words := strings.Split(transcription, " ")

// 	var filteredTimestamps []float64
// 	if len(timestamps) > 0 {
// 		filteredTimestamps = append(filteredTimestamps, timestamps[0]) // Эхний timestamp-г оруулна
// 		for i := 1; i < len(timestamps); i++ {
// 			// Өмнөх timestamp-тай ижил биш бол шинэ массивд нэмнэ
// 			if timestamps[i] != timestamps[i-1] {
// 				filteredTimestamps = append(filteredTimestamps, timestamps[i])
// 			}
// 		}
// 	}

// 	if len(words) > len(timestamps) {
// 		fmt.Println("Анхааруулга: Үгсийн тоо timestamp-үүдээс их байна.")
// 		words = words[:len(timestamps)]
// 	} else if len(words) < len(timestamps) {
// 		fmt.Println("Анхааруулга: Timestamp-үүдийн тоо үгсээс их байна.")
// 		timestamps = timestamps[:len(words)]
// 	}

// 	if err := createSRTFile(words, filteredTimestamps); err != nil {
// 		fmt.Println("Алдаа:", err)
// 		return
// 	}

// 	outputVideo := "output_with_subtitles.mp4"
// 	if err := addSubtitleToVideo(inputVideo, "subtitles.srt", outputVideo); err != nil {
// 		fmt.Println("Алдаа:", err)
// 	}
// }

// func createSRTFile(words []string, filteredTimestamps []float64) error {
// 	file, err := os.Create("subtitles.srt")
// 	if err != nil {
// 		return fmt.Errorf("файлыг үүсгэхэд алдаа: %w", err)
// 	}
// 	defer file.Close()

// 	for i := 0; i < len(words) && i < len(filteredTimestamps); i++ {
// 		startTime := filteredTimestamps[i]
// 		endTime := startTime + 3.0

// 		if i < len(filteredTimestamps)-1 {
// 			endTime = filteredTimestamps[i+1]
// 		}

// 		start := formatTime(startTime)
// 		end := formatTime(endTime)

// 		if _, err := fmt.Fprintf(file, "%d\n%s --> %s\n%s\n\n", i+1, start, end, words[i]); err != nil {
// 			return fmt.Errorf("файлд бичих үед алдаа гарлаа: %w", err)
// 		}
// 	}
// 	fmt.Println("Хадмал текстийн файл 'subtitles.srt' амжилттай үүсгэгдлээ!")
// 	return nil
// }

// func addSubtitleToVideo(inputVideo, subtitleFile, outputVideo string) error {
// 	// Command to add subtitle using FFmpeg
// 	cmd := exec.Command("ffmpeg", "-i", inputVideo, "-vf", fmt.Sprintf("subtitles=%s", subtitleFile), outputVideo)

// 	// Run the FFmpeg command and capture the output
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to add subtitle to video: %w; output: %s", err, output)
// 	}

// 	fmt.Println("Subtitled video created successfully at", outputVideo)
// 	return nil
// }

// func formatTime(seconds float64) string {
// 	hours := int(seconds) / 3600
// 	minutes := (int(seconds) / 60) % 60
// 	secondsInt := int(seconds) % 60
// 	milliseconds := int((seconds - float64(int(seconds))) * 1000)
// 	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secondsInt, milliseconds)
// }
