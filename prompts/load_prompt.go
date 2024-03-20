package prompts

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)
func TemplateExecute(prompt string, data interface{})(string, error){
	templateContent, err := LoadPrompt(prompt)
	if err != nil{
		return "" ,  errors.New(fmt.Sprintln("Error Loading the template:", err))
	}
	// Parse the template
	tmpl, err := template.New("prompt").Parse(templateContent)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Error Parsing the template:", err))
	}
	// Create a buffer to hold the template output
	var buf bytes.Buffer
	// Execute the template with the data
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Error executing template:", err))
	}
	// Get the template output as a string
	return buf.String(), nil
}
func LoadPrompt(prompt string) (string, error) {
	// Assuming voyager is the package name and the prompts directory is in the same directory as the Go file
	packagePath := "vision_builder" // Update this path if necessary

	// Get the absolute path to the prompts directory
	promptsDir := filepath.Join(packagePath, "prompts")

	// Construct the absolute path to the prompt file
	promptFilePath := filepath.Join(promptsDir, prompt+".txt")

	// Read the content of the prompt file
	content, err := os.ReadFile(promptFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to load prompt file: %v", err)
	}

	return string(content), nil
}


