// handlers/chat.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"ai-agent-app/services" // Import the services package

	"github.com/gorilla/mux"
)

// ChatResponse represents the structure of the chat response
type ChatResponse struct {
	Message string `json:"message"`
}

// Global chat history for web requests
var webChatHistory = services.NewChatHistory(10)

// ChatWithAgent handles chat requests with the agent
func ChatWithAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentIDStr := vars["agentID"]

	// Validate and convert agentID to integer
	if agentIDStr == "" {
		http.Error(w, "agentID is required", http.StatusBadRequest)
		return
	}
	
	agentID, err := strconv.Atoi(agentIDStr)
	if err != nil {
		http.Error(w, "Invalid agent ID", http.StatusBadRequest)
		return
	}

	// Extract the message from the request body
	var requestBody struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Validate the message
	if requestBody.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// Add user message to history
	webChatHistory.AddMessage(agentID, "user", requestBody.Message)

	// Get the conversation history
	history := webChatHistory.GetHistory(agentID)

	// Use the OpenAI API to generate a response
	responseMessage, err := services.SendMessageToOpenAI(
		os.Getenv("OPENAI_API_KEY"),
		requestBody.Message,
		history,
	)

	if err != nil {
		http.Error(w, "Error communicating with the agent", http.StatusInternalServerError)
		log.Printf("Error communicating with agent %d: %v", agentID, err)
		return
	}

	// Add the assistant's response to history
	webChatHistory.AddMessage(agentID, "assistant", responseMessage)

	// Log the chat request
	log.Printf("Chat request for agentID: %d, message: %s", agentID, requestBody.Message)

	// Send response
	response := ChatResponse{
		Message: responseMessage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConsoleChatWithAgent handles chat interactions from the console
func ConsoleChatWithAgent(agentID int, message string, chatHistory *services.ChatHistory) (string, error) {
	// Get the conversation history
	history := chatHistory.GetHistory(agentID)
	
	// Use the OpenAI API to generate a response
	responseMessage, err := services.SendMessageToOpenAI(
		os.Getenv("OPENAI_API_KEY"),
		message,
		history,
	)

	if err != nil {
		return "", fmt.Errorf("error communicating with agent %d: %v", agentID, err)
	}

	// Log the console chat request
	log.Printf("Console chat request for agentID: %d, message: %s", agentID, message)

	return responseMessage, nil
}
