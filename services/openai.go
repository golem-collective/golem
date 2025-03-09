package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIAPIURL is the endpoint for the OpenAI API
const OpenAIAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIRequest represents the structure of a request to the OpenAI API
type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"` // Use the Message type from chat_history.go
}

// OpenAIResponse represents the structure of a response from the OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Define a system prompt that sets the context for your AI agent
const defaultSystemPrompt = `You are an AI assistant specialized in helping users with their tasks. 
You are knowledgeable, helpful, and precise in your responses. 
When users ask questions, provide clear and accurate information.
If you're unsure about something, admit it rather than making assumptions.`

// SendMessageToOpenAI sends a message to the OpenAI API and returns the response
func SendMessageToOpenAI(apiKey, userMessage string, history []Message) (string, error) {
	// Create a new request with the user's message and history
	messages := history

	// Add the current message if it's not already in history
	// This is needed because we might be called with just the history
	if len(messages) == 0 || messages[len(messages)-1].Role != "user" || messages[len(messages)-1].Content != userMessage {
		messages = append(messages, Message{
			Role:    "user",
			Content: userMessage,
		})
		messages = append(messages, Message{
			Role:    "user",
			Content: defaultSystemPrompt,
		})
	}

	requestBody := OpenAIRequest{
		Model:    "gpt-3.5-turbo", // You can change this to use a different model
		Messages: messages,
	}

	// Convert the request to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set the headers
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", body)
	}

	// Parse the response
	var openAIResponse OpenAIResponse
	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Check if we got a valid response
	if len(openAIResponse.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	// Return the response
	return openAIResponse.Choices[0].Message.Content, nil
}

// AddMessage is a helper function to add a message to the history
func AddMessage(agentID, role, content string) {
	// This function would typically store the message in a database
	// For now, we'll just log it
	fmt.Printf("Adding message to history for agent %s: %s: %s\n", agentID, role, content)
}
