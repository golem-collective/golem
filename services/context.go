package services

import (
	"strings"
	"text/template"
	"bytes"
)

// ContextBuilder helps build context strings from templates
type ContextBuilder struct {
	template string
	vars     map[string]interface{}
}

// NewContextBuilder creates a new context builder with the given template
func NewContextBuilder(template string) *ContextBuilder {
	return &ContextBuilder{
		template: template,
		vars:     make(map[string]interface{}),
	}
}

// WithVar adds a variable to the context
func (cb *ContextBuilder) WithVar(key string, value interface{}) *ContextBuilder {
	cb.vars[key] = value
	return cb
}

// WithVars adds multiple variables to the context
func (cb *ContextBuilder) WithVars(vars map[string]interface{}) *ContextBuilder {
	for k, v := range vars {
		cb.vars[k] = v
	}
	return cb
}

// Build processes the template and returns the final context string
func (cb *ContextBuilder) Build() (string, error) {
	// Parse the template
	tmpl, err := template.New("context").Parse(cb.template)
	if err != nil {
		return "", err
	}

	// Execute the template with the variables
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, cb.vars); err != nil {
		return "", err
	}

	// Return the processed string
	return buf.String(), nil
}

// BuildContext replaces template variables {{varName}} with values from the state map
func BuildContext(template string, state map[string]string) string {
	result := template

	// Replace each variable in the template with its value from state
	for key, value := range state {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// Example usage:
// template := `You are {{role}}, an AI assistant specialized in {{specialty}}.
// Your task is to {{task}}.
// Use these key traits in your responses: {{traits}}`
//
// state := map[string]string{
//     "role": "Tech Expert",
//     "specialty": "programming",
//     "task": "help users write better code",
//     "traits": "clear, precise, and educational",
// }
//
// context := BuildContext(template, state) 