// services/agent.go
package services

import (
	"ai-agent-app/database"
	"ai-agent-app/models"
	"log"
)

// CreateAgent saves a new agent to the database and returns its ID
func CreateAgent(agent *models.Agent) error {
	// Prepare the SQL statement with RETURNING clause to get the generated ID
	query := `INSERT INTO agents (name, type, context) VALUES ($1, $2, $3) RETURNING id`
	err := database.GetDB().QueryRow(query, agent.Name, agent.Type, agent.Context).Scan(&agent.ID)
	if err != nil {
		log.Printf("Error saving agent to database: %v", err)
		return err
	}
	return nil
}
