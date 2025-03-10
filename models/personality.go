package models

type Personality struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	Description string   `json:"description"`
	System      string   `json:"system"`
	Bio         []string `json:"bio"`
	Lore        []string `json:"lore"`
	Knowledge   []string `json:"knowledge"`
	Style       struct {
		All  []string `json:"all"`
		Chat []string `json:"chat"`
	} `json:"style"`
	Adjectives   []string `json:"adjectives"`
	Instructions string   `json:"instructions"`
}
