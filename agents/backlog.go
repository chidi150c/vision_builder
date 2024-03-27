package agents

import (
	"ai_agents/vision_builder/prompts"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	// "os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type CurriculumAgent struct {
	// Llm                   *ChatOpenAI // Assume ChatOpenAI is a defined struct for handling AI chat operations
	// QaLlm                 *ChatOpenAI
	Mode                  string
	CkptDir               string
	Resume                bool
	CompletedTasks        []string
	FailedTasks           []string
	QACache               map[string]string
	// QaCacheQuestionsVectordb *VectorDB // Assume VectorDB is a defined struct for vector operations
	WarmUp                map[string]int
	DefaultWarmUp         map[string]int
	CurriculumObservations []string
	CoreInventoryItemsRegex *regexp.Regexp
	mu                    sync.RWMutex // For concurrent access to the QaCache
	AbsolutePath string
}

func NewCurriculumAgent(abPath string, temperature float64, qaModelName string, qaTemperature float64, requestTimeout time.Duration, ckptDir string, resume bool, mode string, warmUp map[string]int, coreInventoryItems string) *CurriculumAgent {
    // reg, _ := regexp.Compile(coreInventoryItems)
    return &CurriculumAgent{
        // ModelName:             modelName,
        // Temperature:           temperature,
        // QAModelName:           qaModelName,
        // QATemperature:         qaTemperature,
        // RequestTimeout:        requestTimeout,
        AbsolutePath: abPath,
        CkptDir:               ckptDir,
        Resume:                resume,
        Mode:                  mode,
        WarmUp:                warmUp,
        // CoreInventoryItemsReg: reg,
        CompletedTasks:        make([]string, 0),
        FailedTasks:           make([]string, 0),
        QACache:               make(map[string]string),
        // QACacheQuestionsVectorDB: NewVectorDB("qaCacheQuestionsVectorDB", ckptDir+"/curriculum/vectordb"), // Assume NewVectorDB is a constructor for VectorDB
        // ChatModel:               NewChatOpenAI(modelName, temperature, requestTimeout), // Assume NewChatOpenAI is a constructor for ChatOpenAI
    }
}

// RenderSystemMessage generates a system message based on a predefined prompt.
func (ca *CurriculumAgent) RenderSystemMessage() string {
    message, err := prompts.LoadPrompt(ca.AbsolutePath, "models")
    if err != nil {
        log.Printf("Error loading system message prompt: %v\n", err)
        return ""
    }
    return message
}

// ProposeNextTask suggests the next task based on the agent's mode and progress.
func (ca *CurriculumAgent) ProposeNextTask() (string, string) {
    if ca.Mode == "auto" {
        // Simplified version: returns fixed task and context based on some conditions
        if len(ca.CompletedTasks) == 0 {
            return "Mine 1 wood log", "You can mine one of oak, birch, spruce, jungle, acacia, dark oak, or mangrove logs."
        }
        // Additional logic to determine the next task based on completed tasks and mode...
    } else if ca.Mode == "manual" {
        // Handle manual task proposal, potentially involving user input or another mechanism
    }
    return "", "" // Fallback if no task is proposed
}


// LoadJson loads JSON data from a file into a given struct.
func LoadJson(filePath string, v interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &v)
}

// DumpJson writes the given struct as JSON to a file.
func DumpJson(v interface{}, filePath string) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0644)
}
// SaveState writes the agent's state to JSON files.
func (ca *CurriculumAgent) SaveState() error {
    if err := ca.saveToFile(ca.CompletedTasks, filepath.Join(ca.CkptDir, "curriculum", "completed_tasks.json")); err != nil {
        return err
    }
    if err := ca.saveToFile(ca.FailedTasks, filepath.Join(ca.CkptDir, "curriculum", "failed_tasks.json")); err != nil {
        return err
    }
    if err := ca.saveToFile(ca.QACache, filepath.Join(ca.CkptDir, "curriculum", "qa_cache.json")); err != nil {
        return err
    }
    return nil
}

// saveToFile serializes the given data to a JSON file.
func (ca *CurriculumAgent) saveToFile(data interface{}, filepath string) error {
    bytes, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filepath, bytes, 0644)
}

// GenerateQuestions creates a list of questions based on the agent's current state and observations.
func (ca *CurriculumAgent) GenerateQuestions(observation string) ([]string, error) {
    // This method would dynamically generate questions based on the current
    // state of the environment and what the agent has observed so far.
    // The actual implementation will depend on how observations are structured
    // and what information is relevant for generating useful questions.

    // Example placeholder logic:
    if strings.Contains(observation,"tree") {
        return []string{"What can be made from wood in Minecraft?"}, nil
    }

    return nil, fmt.Errorf("question generation not implemented")
}
