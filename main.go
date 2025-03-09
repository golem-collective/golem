package main

import (
	"ai-agent-app/database"
	"ai-agent-app/handlers"
	"ai-agent-app/services"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func main() {
	fmt.Println("AI Agent Console")
	fmt.Println("----------------")

	// Initialize database connection
	database.InitDB()
	defer database.CloseDB()

	// Create necessary tables
	if err := database.CreateAgentsTable(); err != nil {
		log.Fatalf("Failed to create agents table: %v", err)
	}
	if err := database.CreateChatHistoryTable(); err != nil {
		log.Fatalf("Failed to create chat history table: %v", err)
	}

	// For debugging - print the API key (remove in production)
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: OPENAI_API_KEY is not set")
	} else {
		log.Println("OPENAI_API_KEY is set")
	}

	// Initialize chat history service
	chatHistory := services.NewChatHistory(10) // Keep last 10 messages

	// Create a default agent
	agentID, err := handlers.CreateDefaultAgent()
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	fmt.Printf("Agent created with ID: %d\n", agentID)
	fmt.Println("Start chatting with the agent (type 'exit' to quit, 'clear' to clear history):")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if strings.ToLower(userInput) == "clear" {
			chatHistory.ClearHistory(agentID)
			fmt.Println("Chat history cleared.")
			continue
		}

		// Add user message to history
		chatHistory.AddMessage(agentID, "user", userInput)

		// Chat with the agent
		response, err := handlers.ConsoleChatWithAgent(agentID, userInput, chatHistory)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Add agent response to history
		chatHistory.AddMessage(agentID, "assistant", response)

		fmt.Printf("Agent: %s\n", response)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
