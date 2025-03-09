package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// GrokAPIURL is the endpoint for the Grok API
const GrokAPIURL = "https://api.grok.com/v1/chat"

// GrokResponse represents the response structure from the Grok API
type GrokResponse struct {
	Response string `json:"response"`
}

// SendMessageToGrok sends a message to the Grok API and returns the response
func SendMessageToGrok(message string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"message": message,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(GrokAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var grokResponse GrokResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResponse); err != nil {
		return "", err
	}

	return grokResponse.Response, nil
}
