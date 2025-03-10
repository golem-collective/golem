package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ai-agent-app/models"
)

// LoadPersonality loads a personality from a JSON file
func LoadPersonality(agentName string) (*models.Personality, error) {
	// Construct the file path
	filePath := filepath.Join("personalities", fmt.Sprintf("%s.json", agentName))

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read personality file: %w", err)
	}

	// Parse the JSON data
	var personality models.Personality
	if err := json.Unmarshal(data, &personality); err != nil {
		return nil, fmt.Errorf("failed to parse personality data: %w", err)
	}

	return &personality, nil
}
