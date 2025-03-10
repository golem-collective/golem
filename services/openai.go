package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// OpenAIAPIURL is the endpoint for the OpenAI API
const OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIRequest represents the structure of a request to the OpenAI API
type OpenAIRequest struct {
	Model   string `json:"model"`
	Context string `json:"context"`
}

// OpenAIResponse represents the structure of a response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// SendMessageToOpenAI sends a message to the OpenAI API and returns the response
func SendMessageToOpenAI(apiKey string, userMessage string, systemTemplate string, state map[string]string) (string, error) {
	// Replace template variables with values from state
	systemMessage := systemTemplate
	for key, value := range state {
		systemMessage = strings.Replace(systemMessage, "{{"+key+"}}", value, -1)
	}

	var context = BuildContext(systemTemplate, state)

	// Create the request body
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": userMessage,
			},
			{
				"role":    "system",
				"content": context,
			},
		},
	}

	requestData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", string(body))
	}

	// Parse the response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Extract the message content
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid choice format")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format")
	}

	return content, nil
}

// AddMessage is a helper function to add a message to the history
func AddMessage(agentID, role, content string) {
	// This function would typically store the message in a database
	// For now, we'll just log it
	fmt.Printf("Adding message to history for agent %s: %s: %s\n", agentID, role, content)
}
