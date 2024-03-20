package config

import "os"


type ModelConfig struct{	
    UserInput string
	ApiKey string
	Url string
	Model string
}

func NewModelConfigs() map[string]ModelConfig {
	return map[string]ModelConfig{
		"gpt4": {
			UserInput : "your_topic_here", // Replace with your actual input
			ApiKey : os.Getenv("OPENAI_API_KEY"), // Ensure you have set your API key in your environment variables
			Url : "https://api.openai.com/v1/chat/completions",
			Model:    "gpt-3.5-turbo",  //"gpt-3.5-turbo", //"gpt-4",
		},
		"ollama": {
			Url: "http://localhost:11434/api/generate",
			Model:     "mistral", 
		},
	}
}
