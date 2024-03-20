package builder

import (
	"ai_agents/vision_builder/agents"
	"ai_agents/vision_builder/config"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)


type VisionBuilder struct {
	ID              int
	VisionStatement           string
	ActionAgentRolloutNumIter int
	ActionAgent               *agents.ActionAgent
	AI agents.CloseAIModelServicer
	Messages []agents.Message
	Conversations []string
	Output          CodeResponse
	AllBuilder      strings.Builder
	UnMappedBackLog []string
	ResearchBackLog []string
	Reader          *bufio.Reader
	Fnum            int
	FilePath        string
	Vision          string
	Goal            string
	Task            string
	SubTask         string
}

func NewVisionBuilder(visionStatement string) *VisionBuilder {
	// Define AI configuration
	config := config.NewModelConfigs()["gpt4"]
	// Create an instance of the real AI model
	openAI := agents.NewOpenAI(config)	
	closeAI := agents.NewCloseAIModel(config)
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("error Getting the current working directory")
	}
	// Construct the package path relative to the current working directory
	ckptDir := filepath.Join(currentDir, "vision_builder")
	AIAgent := agents.NewActionAgent(openAI, ckptDir, false, false, false)
	return &VisionBuilder{
		VisionStatement: visionStatement,
		ActionAgent:     AIAgent,
		AI: closeAI,
	}
}

// Step progresses the state of the task execution by one step.
func (v *VisionBuilder) Step() ([]agents.Message, int, bool, map[string]interface{}) {
	if v.ActionAgentRolloutNumIter < 0 {
		panic("Agent must be reset before stepping")
	}
	//Process messages through an AI model (LLM) and possibly executes generated code in an environment.
	// []Message{
	// 	{"system", msg[0]},
	// 	{"user", msg[1]},
	// }
	aiMessage, err := v.ActionAgent.Llm.PromptAI(v.Messages)
	if err != nil{
		log.Fatalln(err)
	}
	fmt.Printf("\033[34m****Action Agent ai message****\n%s\033[0m\n", aiMessage.Content)
	v.Conversations = append(v.Conversations, v.Messages[0].Content)
	v.Conversations = append(v.Conversations, v.Messages[1].Content)
	v.Conversations = append(v.Conversations, aiMessage.Content)
	parsedResult, err := v.ActionAgent.ProcessAIMessage(aiMessage)
	if err != nil{
		log.Fatalln(err)
	}
	// var success bool
	// var critique string
	// var code string

	fmt.Println(parsedResult)
	panic("My testing see above parseResult")
	// if parsedResult != nil {
	// 	switch result := parsedResult.(type) {
	// 	case map[string]interface{}:
	// 		code = result["program_code"].(string) + "\n" + result["exec_code"].(string)
	// 		events := v.Env.Step(code, v.SkillManager.Programs)
	// 		v.Recorder.Record(events, v.Task)
	// 		v.ActionAgent.UpdateChestMemory(events[len(events)-1][1]["nearbyChests"])
	// 		success, critique = v.CriticAgent.CheckTaskSuccess(events, v.Task, v.Context, v.ActionAgent.RenderChestObservation(), 5)

	// 		if v.ResetPlacedIfFailed && !success {
	// 			// Revert all placing events in the last step
	// 			var blocks []string
	// 			var positions []interface{}
	// 			for _, event := range events {
	// 				if event[0] == "onSave" && strings.HasSuffix(event[1]["onSave"].(string), "_placed") {
	// 					block := strings.TrimSuffix(event[1]["onSave"].(string), "_placed")
	// 					position := event[1]["status"].(map[string]interface{})["position"]
	// 					blocks = append(blocks, block)
	// 					positions = append(positions, position)
	// 				}
	// 			}
	// 			newEvents := v.Env.Step(fmt.Sprintf("await givePlacedItemBack(bot, %s, %s)", json.dumps(blocks), json.dumps(positions)), v.SkillManager.Programs)
	// 			events[len(events)-1][1]["inventory"] = newEvents[len(newEvents)-1][1]["inventory"]
	// 			events[len(events)-1][1]["voxels"] = newEvents[len(newEvents)-1][1]["voxels"]
	// 		}
	// 		newSkills := v.SkillManager.RetrieveSkills(v.Context + "\n\n" + v.ActionAgent.SummarizeChatlog(events))
	// 		systemMessage := v.ActionAgent.RenderSystemMessage(newSkills)
	// 		humanMessage := v.ActionAgent.RenderHumanMessage(events, code, v.Task, v.Context, critique)
	// 		v.LastEvents = events
	// 		v.Messages = []string{systemMessage, humanMessage}
	// 	case string:
	// 		v.Recorder.Record([][]map[string]interface{}{}, v.Task)
	// 		fmt.Printf("\033[34m%s Trying again!\033[0m\n", result)
	// 	}
	// } else {
	// 	panic("parsed_result is nil")
	// }
	// v.ActionAgentRolloutNumIter++
	// done := v.ActionAgentRolloutNumIter >= v.ActionAgentTaskMaxRetries || success
	// info := map[string]interface{}{
	// 	"task":          v.Task,
	// 	"success":       success,
	// 	"conversations": v.Conversations,
	// }
	// if success {
	// 	info["program_code"] = parsedResult["program_code"]
	// 	info["program_name"] = parsedResult["program_name"]
	// } else {
	// 	fmt.Printf("\033[32m****Action Agent human message****\n%s\033[0m\n", v.Messages[len(v.Messages)-1])
	// }
	// return v.Messages, 0, done, info
}  

func (v *VisionBuilder) Rollout(task string, context string, resetEnv bool) ([]agents.Message, int, bool, map[string]interface{}) {
	// v.Reset(task, context, resetEnv)
	// for {
	// 	messages, reward, done, info := v.Step()
	// 	if done {
	// 		break
	// 	}
	// }
	return v.Messages, 0, false, nil
}