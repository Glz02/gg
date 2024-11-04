// audio.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func pushAudio(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	req, err := http.NewRequest("POST", "https://api.chimege.com/v1.2/stt-long-hq", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("token", "459af72480c5cdc493187acfcb04bac8813227e832c727b4c4d633f07fe14698")
	req.Header.Set("Mode", "offsets")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var res Response
	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	return res.UUID, nil
}

func getOutputText(uuid string) (*TranscribedText, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.chimege.com/v1.2/stt-long-hq-transcript"),
		nil,
	)

	req.Header.Set("Token", "459af72480c5cdc493187acfcb04bac8813227e832c727b4c4d633f07fe14698")
	req.Header.Set("UUID", uuid)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := new(TranscribedText)
	if resp.StatusCode != 200 {
		return nil, errors.New("Sorry unknown error occurred")
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
