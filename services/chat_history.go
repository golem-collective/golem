package services

import (
	"ai-agent-app/database"
	"log"
)

// Message represents a single message in the chat history
type Message struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // The message content
}

// ChatHistory stores conversation history for each agent
type ChatHistory struct {
	maxLength int
}

// NewChatHistory creates a new chat history manager
func NewChatHistory(maxLength int) *ChatHistory {
	return &ChatHistory{
		maxLength: maxLength,
	}
}

// AddMessage adds a message to the conversation history for a specific agent
func (ch *ChatHistory) AddMessage(agentID int, role, content string) {
	query := `
		INSERT INTO chat_history (agent_id, role, content)
		VALUES ($1, $2, $3)`
	
	_, err := database.Exec(query, agentID, role, content)
	if err != nil {
		log.Printf("Error adding message to chat history: %v", err)
		return
	}

	// Trim history if it exceeds maximum length
	trimQuery := `
		DELETE FROM chat_history 
		WHERE id IN (
			SELECT id FROM chat_history 
			WHERE agent_id = $1 
			ORDER BY created_at DESC 
			OFFSET $2
		)`
	
	_, err = database.Exec(trimQuery, agentID, ch.maxLength)
	if err != nil {
		log.Printf("Error trimming chat history: %v", err)
	}
}

// GetHistory returns the conversation history for a specific agent
func (ch *ChatHistory) GetHistory(agentID int) []Message {
	query := `
		SELECT role, content 
		FROM chat_history 
		WHERE agent_id = $1 
		ORDER BY created_at ASC`

	db := database.GetDB()
	rows, err := db.Query(query, agentID)
	if err != nil {
		log.Printf("Error getting chat history: %v", err)
		return []Message{}
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.Role, &msg.Content); err != nil {
			log.Printf("Error scanning chat history row: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating chat history rows: %v", err)
	}

	return messages
}

// ClearHistory clears the conversation history for a specific agent
func (ch *ChatHistory) ClearHistory(agentID int) {
	query := `DELETE FROM chat_history WHERE agent_id = $1`
	_, err := database.Exec(query, agentID)
	if err != nil {
		log.Printf("Error clearing chat history: %v", err)
	}
}
