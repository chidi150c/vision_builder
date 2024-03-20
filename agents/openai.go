package agents

import (
	"ai_agents/vision_builder/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Define structs to match the response structure
type OpenAIResponse struct {
	ID       string  `json:"id"`
	Object   string  `json:"object"`
	Created  int64   `json:"created"`
	Model    string  `json:"model"`
	Choices  []Choice `json:"choices"`
	Usage    Usage   `json:"usage"`
}

type Choice struct {
	Index         int    `json:"index"`
	Message       Message `json:"message"`
	Logprobs      interface{} `json:"logprobs"`
	FinishReason  string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens    int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens     int `json:"total_tokens"`
}

// Define a struct to hold the request payload
type CompletionRequest struct {
	Model    string `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIModel struct{
	Name string
	ApiKey string
	Url string
	Model string 
}

func NewOpenAI(cf config.ModelConfig)*OpenAIModel{
	return &OpenAIModel{
		ApiKey: cf.ApiKey,
		Url: cf.Url, 
		Model: cf.Model,
	}
}

// AIModel represents an interface for interacting with an AI model.
type AIModelServicer interface {
	PromptAI(messages []Message) (Message, error)
}

var _ AIModelServicer = &OpenAIModel{}

func apiFetch(messages []Message, Model, Url, ApiKey string)(Message, error){
	requestBody := CompletionRequest{
		Model:    Model,
		Messages: messages,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return Message{}, fmt.Errorf("error marshalling the request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return Message{}, fmt.Errorf("error creating the request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("error making the request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Message{}, fmt.Errorf("error reading the response body: %v", err)
	}
	var response OpenAIResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return Message{}, fmt.Errorf("error parsing the JSON response: %v", err)
	}

	// Assuming you're interested in the content of the first choice
	if len(response.Choices) > 0 {
		return response.Choices[0].Message, nil
	} else {
		return Message{}, fmt.Errorf("no content generated")
	}
}
func (a *OpenAIModel)PromptAI(messages []Message) (Message, error){
	content, err := apiFetch(messages, a.Model,a.Url, a.ApiKey)
	if err != nil {
		log.Fatalf("error fetching API content: %v", err)
	}		
	return content, nil
}



type CloseAIModel struct{
	Name string
	ApiKey string
	Url string
	Model string 
}
func NewCloseAIModel(cf config.ModelConfig)*CloseAIModel{
	return &CloseAIModel{
		ApiKey: cf.ApiKey,
		Url: cf.Url, 
		Model: cf.Model,
	}
}
// AIModel represents an interface for interacting with an AI model.
type CloseAIModelServicer interface {
	PromptAI(system, user string) (string, error)
}
var _ CloseAIModelServicer = &CloseAIModel{}
func (a *CloseAIModel)PromptAI(system, user string) (string, error){
	var messages = []Message{
		{"system", system},
		{"user", user},
	}
	content, err := apiFetchClose(messages, a.Model,a.Url, a.ApiKey)
	if err != nil {
		log.Fatalf("error fetching API content: %v", err)
	}		
	return content, nil
}
func apiFetchClose(messages []Message, Model, Url, ApiKey string)(string, error){
	requestBody := CompletionRequest{
		Model:    Model,
		Messages: messages,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshalling the request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", Url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating the request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making the request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %v", err)
	}
	var response OpenAIResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("error parsing the JSON response: %v", err)
	}

	// Assuming you're interested in the content of the first choice
	if len(response.Choices) > 0 {
		generatedContent := response.Choices[0].Message.Content
		return generatedContent, nil
	} else {
		return "", fmt.Errorf("no content generated")
	}
}