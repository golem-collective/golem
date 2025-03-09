package models

type Agent struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"` // e.g., "openai" or "grok"
	Context     string `json:"context"`
}
