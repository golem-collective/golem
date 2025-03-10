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
	query := `INSERT INTO agents (name) VALUES ($1) RETURNING id`
	err := database.GetDB().QueryRow(query,
		agent.Name,
	).Scan(&agent.ID)

	if err != nil {
		log.Printf("Error saving agent to database: %v", err)
		return err
	}
	return nil
}

// GetAgentByID retrieves an agent from the database by its ID
func GetAgentByID(id int) (*models.Agent, error) {
	query := `SELECT id, name FROM agents WHERE id = $1`

	var agent models.Agent
	err := database.GetDB().QueryRow(query, id).Scan(
		&agent.ID,
		&agent.Name,
	)

	if err != nil {
		log.Printf("Error retrieving agent with ID %d: %v", id, err)
		return nil, err
	}

	return &agent, nil
}
