// handlers/chat.go
package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"ai-agent-app/services" // Import the services package
)

// ChatResponse represents the structure of the chat response
type ChatResponse struct {
	Message string `json:"message"`
}

// Global chat history for web requests
// var webChatHistory = services.NewChatHistory(10)

// ChatWithAgent handles chat requests with the agent
// func ChatWithAgent(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	agentIDStr := vars["agentID"]

// 	// Validate and convert agentID to integer
// 	if agentIDStr == "" {
// 		http.Error(w, "agentID is required", http.StatusBadRequest)
// 		return
// 	}

// 	agentID, err := strconv.Atoi(agentIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid agent ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Extract the message from the request body
// 	var requestBody struct {
// 		Message string `json:"message"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		log.Printf("Error decoding request body: %v", err)
// 		return
// 	}

// 	// Validate the message
// 	if requestBody.Message == "" {
// 		http.Error(w, "Message is required", http.StatusBadRequest)
// 		return
// 	}

// 	// Add user message to history
// 	webChatHistory.AddMessage(agentID, "user", requestBody.Message)

// 	// Get the conversation history
// 	history := webChatHistory.GetHistory(agentID)

// 	// Define the context template with conversation history
// 	template := `You are {{name}}, an AI assistant.
// {{description}}
// You specialize in {{specialty}}.
// Communication style: {{style}}
// Follow these rules:
// {{rules}}

// Previous conversation:
// {{history}}

// Current request:
// `

// 	agent, err := services.GetAgentByID(agentID)
// 	// Create state map for context variables
// 	state := map[string]string{
// 		"name":        agent.Name,
// 		"description": agent.Description,
// 		"specialty":   agent.System,
// 		"style":       "professional and helpful",
// 		"rules":       "1. Be concise\n2. Be accurate\n3. Be helpful",
// 		"history":     formatHistory(history), // You'll need to implement this
// 	}

// 	// Use the OpenAI API to generate a response
// 	responseMessage, err := services.SendMessageToOpenAI(
// 		os.Getenv("OPENAI_API_KEY"),
// 		requestBody.Message,
// 		template,
// 		state,
// 	)

// 	if err != nil {
// 		http.Error(w, "Error communicating with the agent", http.StatusInternalServerError)
// 		log.Printf("Error communicating with agent %d: %v", agentID, err)
// 		return
// 	}

// 	// Add the assistant's response to history
// 	webChatHistory.AddMessage(agentID, "assistant", responseMessage)

// 	// Log the chat request
// 	log.Printf("Chat request for agentID: %d, message: %s", agentID, requestBody.Message)

// 	// Send response
// 	response := ChatResponse{
// 		Message: responseMessage,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// ConsoleChatWithAgent handles chat interactions from the console
func ConsoleChatWithAgent(agentID int, message string, chatHistory *services.ChatHistory) (string, error) {
	// Get the conversation history
	history := chatHistory.GetHistory(agentID)

	agent, err := services.GetAgentByID(agentID)

	if err != nil {
		return "", fmt.Errorf("error retrieving agent %d: %v", agentID, err)
	}

	personality, err := services.LoadPersonality(agent.Name)

	if err != nil {
		return "", fmt.Errorf("error loading personality: %w", err)
	}

	template := `You are {{name}}, an AI assistant.
		{{description}}
		{{system}}
		
		Background:
		{{bio}}
		
		Experience:
		{{lore}}
		
		Expertise:
		{{knowledge}}
		
		Communication style:
		{{style}}
		
		Instructions:
		{{instructions}}
		
		Chat history:
		{{history}}
		
		Adjectives:
		{{adjectives}}
		
		Instructions:
		{{instructions}}
		
		`

	// Create state map for context variables
	state := map[string]string{
		"name":         agent.Name,
		"description":  personality.Description,
		"specialty":    personality.System,
		"history":      formatChatHistory(history),
		"style":        strings.Join(personality.Style.Chat, "\n"),
		"bio":          strings.Join(personality.Bio, "\n"),
		"lore":         strings.Join(personality.Lore, "\n"),
		"knowledge":    strings.Join(personality.Knowledge, "\n"),
		"adjectives":   strings.Join(personality.Adjectives, "\n"),
		"instructions": personality.Instructions,
	}

	// Use the OpenAI API to generate a response
	responseMessage, err := services.SendMessageToOpenAI(
		os.Getenv("OPENAI_API_KEY"),
		message,
		template,
		state,
	)

	if err != nil {
		return "", fmt.Errorf("error communicating with agent %d: %v", agentID, err)
	}

	// Log the console chat request
	log.Printf("Console chat request for agentID: %d, message: %s", agentID, message)

	return responseMessage, nil
}

// formatHistory converts the chat history array to a formatted string
func formatHistory(history []services.Message) string {
	if len(history) == 0 {
		return "No previous conversation."
	}

	var formattedHistory string
	for _, msg := range history {
		formattedHistory += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}

	return formattedHistory
}

// formatChatHistory is an alias for formatHistory for consistency
func formatChatHistory(history []services.Message) string {
	return formatHistory(history)
}
