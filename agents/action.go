package agents

import (
	// "context"
	"ai_agents/vision_builder/prompts"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// "log"
	"regexp"
	"strings"
)

type ActionAgent struct {
	CkptDir        string
	ChatLog        bool
	ExecutionError bool
	AppMemory    map[string]interface{}
	Llm            AIModelServicer
}

func NewActionAgent(AIModel AIModelServicer,ckptDir string, resume bool, chatLog bool, executionError bool) *ActionAgent {
	agent := &ActionAgent{
		CkptDir:        ckptDir,
		ChatLog:        chatLog,
		ExecutionError: executionError,
		AppMemory:    make(map[string]interface{}),
		Llm: AIModel,
	}
		
	// Create a new scanner to read input
	scanner := bufio.NewScanner(os.Stdin)

	// Ask the user for the directory name
	fmt.Printf("Enter a directory name where the app would be set up (defualt: \"%s\"): ", ckptDir)
	// Read the input from the user
	if scanner.Scan() {
		// Get the directory name entered by the user
		dirName := strings.TrimSpace(scanner.Text())
		// Check if the directory name is empty
		if dirName != "" {
			ckptDir = dirName
		}
		agent.CkptDir = ckptDir
		// Create the directory if it doesn't exist
		err := os.MkdirAll(fmt.Sprintf("%s/action", ckptDir), 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			log.Fatalln(err)
		}
		fmt.Println("Directory created successfully")
	} else {
		fmt.Println("Failed to read input.")
	}
	// Specify the path to the JSON file
	filePath := fmt.Sprintf("%s/action/app_memory.json", ckptDir)

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// Initialize a variable to store the decoded JSON data
	var appMemory map[string]interface{}

	// Unmarshal the JSON data into the appMemory variable
	err = json.Unmarshal(data, &appMemory)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	// Access the loaded JSON data
	// For example:
	fmt.Println("Loaded app memory:")
	for position, app := range appMemory {
		if appMap, ok := app.(map[string]interface{}); ok{
			fmt.Printf("Position: %s, App: %+v\n", position, appMap)
			if resume {
				fmt.Printf("\033[32mLoading Action Agent from %s/action\033[0m\n", ckptDir)
				// appMemory := utils.LoadJSON(fmt.Sprintf("%s/action/app_memory.json", ckptDir))
				agent.AppMemory[position] = appMap
			}
		}
	}
	return agent
}

// DumpJSON dumps the AppMemory data to a JSON file
func (agent *ActionAgent) DumpJSON() error {
	// Marshal the AppMemory data into JSON format
	jsonData, err := json.MarshalIndent(agent.AppMemory, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON data: %v", err)
	}

	// Specify the path to the JSON file
	filePath := fmt.Sprintf("%s/action/app_memory.json", agent.CkptDir)

	// Write the JSON data to the file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	fmt.Println("JSON data successfully written to file:", filePath)
	return nil
}

func (agent *ActionAgent) UpdateAppMemory(apps map[string]interface{}) {
	for position, app := range apps {
		if existingApp, ok := agent.AppMemory[position]; ok {
			if _, ok := existingApp.(map[string]interface{}); ok {
				if appMap, ok := app.(map[string]interface{}); ok {
					agent.AppMemory[position] = appMap
				}
			}
			if app == "Invalid" {
				fmt.Printf("\033[32mAction Agent removing app %s: %v\033[0m\n", position, app)
				delete(agent.AppMemory, position)
			}
		} else {
			if app != "Invalid" {
				fmt.Printf("\033[32mAction Agent saving app %s: %v\033[0m\n", position, app)
				agent.AppMemory[position] = app
			}
		}
	}
	agent.DumpJSON()
}

func (agent *ActionAgent) RenderAppObservation() string {
	var observations []string
	for appPosition, app := range agent.AppMemory {
		if appMap, ok := app.(map[string]interface{}); ok && len(appMap) > 0 {
			observations = append(observations, fmt.Sprintf("%s: %v", appPosition, app))
		} else if appMap, ok := app.(map[string]interface{}); ok && len(appMap) == 0 {
			observations = append(observations, fmt.Sprintf("%s: Empty", appPosition))
		} else if appStr, ok := app.(string); ok && appStr == "Unknown" {
			observations = append(observations, fmt.Sprintf("%s: Unknown items inside", appPosition))
		}
	}
	if len(observations) > 0 {
		return "Apps:\n" + strings.Join(observations, "\n") + "\n\n"
	}
	return "Apps: None\n\n"
}
type SystemTemplate struct{
	Programs string
	ResponseFormat string
}
func (agent *ActionAgent) RenderSystemMessage(tasks []string) string {
	baseTasks := []string{
		"exploreUntil",
		"mineBlock",
		"craftItem",
		"placeItem",
		"smeltItem",
		"killMob",
	}

	programs := strings.Join(agent.LoadControlPrimitivesContext(baseTasks), "\n") + "\n" + strings.Join(tasks, "\n")
	responseFormat, err := prompts.LoadPrompt("action_response_format")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	sm := SystemTemplate{
		Programs: programs,
		ResponseFormat: responseFormat,
	}
	systemMessagePrompt, err := prompts.TemplateExecute("action_template", sm)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	return systemMessagePrompt
}

func (agent *ActionAgent) RenderHumanMessage(events []interface{}, code string, task string, context string, critique string) string {
	var chatMessages []string
	var errorMessages []string
	var damageMessages []string
	var inventoryUsed int
	var inventory interface{}

	_ = func(cond bool, msg string) {
		if !cond {
			panic(msg)
		}
	}

	for _, event := range events {
		switch event.(type) {
		case map[string]interface{}:
			eventMap := event.(map[string]interface{})
			eventType := eventMap["eventType"].(string)
			switch eventType {
			case "onChat":
				chatMessages = append(chatMessages, eventMap["onChat"].(string))
			case "onError":
				errorMessages = append(errorMessages, eventMap["onError"].(string))
			case "onDamage":
				damageMessages = append(damageMessages, eventMap["onDamage"].(string))
			case "observe":
				status := eventMap["status"].(map[string]interface{})
				inventoryUsed = status["inventoryUsed"].(int)
				inventory = eventMap["inventory"]
			}
		}
	}

	observation := ""
	if code != "" {
		observation += fmt.Sprintf("Code from the last round:\n%s\n\n", code)
	} else {
		observation += "Code from the last round: No code in the first round\n\n"
	}

	if agent.ExecutionError {
		if len(errorMessages) > 0 {
			observation += fmt.Sprintf("Execution error:\n%s\n\n", strings.Join(errorMessages, "\n"))
		} else {
			observation += "Execution error: No error\n\n"
		}
	}

	if agent.ChatLog {
		if len(chatMessages) > 0 {
			observation += fmt.Sprintf("Chat log: %s\n\n", strings.Join(chatMessages, "\n"))
		} else {
			observation += "Chat log: None\n\n"
		}
	}

	if inventory != nil {
		observation += fmt.Sprintf("Inventory (%d/36): %v\n\n", inventoryUsed, inventory)
	} else {
		observation += "Inventory (0/36): Empty\n\n"
	}

	observation += agent.RenderAppObservation()
	observation += fmt.Sprintf("Task: %s\n\n", task)
	observation += fmt.Sprintf("Context: %s\n\n", context)
	observation += fmt.Sprintf("Critique: %s\n\n", critique)

	return observation
}

func (agent *ActionAgent) ProcessAIMessage(message Message) (string, error) {
	var errorStr string
	codePattern := regexp.MustCompile("```(?:golang|go)(.*?)```")
	codeMatches := codePattern.FindAllStringSubmatch(message.Content, -1)
	var codeLines []string
	for _, match := range codeMatches {
		codeLines = append(codeLines, match[1])
	}
	code := strings.Join(codeLines, "\n")
	return code, fmt.Errorf(errorStr)
}


func (agent *ActionAgent) SummarizeChatlog(events []interface{}) string {
	filterItem := func(message string) string {
		craftPattern := regexp.MustCompile(`I cannot make \w+ because I need: (.*)`)
		craftPattern2 := regexp.MustCompile(`I cannot make \w+ because there is no crafting table nearby`)
		minePattern := regexp.MustCompile(`I need at least a (.*) to mine \w+!`)
		if craftPattern.MatchString(message) {
			return craftPattern.FindStringSubmatch(message)[1]
		} else if craftPattern2.MatchString(message) {
			return "a nearby crafting table"
		} else if minePattern.MatchString(message) {
			return minePattern.FindStringSubmatch(message)[1]
		}
		return ""
	}

	chatlog := make(map[string]bool)
	for _, event := range events {
		eventMap := event.(map[string]interface{})
		eventType := eventMap["eventType"].(string)
		if eventType == "onChat" {
			chatMessage := eventMap["onChat"].(string)
			item := filterItem(chatMessage)
			if item != "" {
				chatlog[item] = true
			}
		}
	}

	var items []string
	for item := range chatlog {
		items = append(items, item)
	}

	if len(items) > 0 {
		return "I also need " + strings.Join(items, ", ") + "."
	}
	return ""
}

func (agent *ActionAgent) LoadControlPrimitivesContext(primitiveNames []string)[]string {
	// Usage example
	// primitiveNames Pass primitive names if available, otherwise it will be auto-detected
	primitives, err := loadControlPrimitivesContext(primitiveNames, agent.CkptDir)
	if err != nil {
		log.Fatalln("Error:", err)
		return nil
	}
	fmt.Println("Loaded primitives:")
	for _, primitive := range primitives {
		fmt.Println(primitive)
	}
	return primitives
}