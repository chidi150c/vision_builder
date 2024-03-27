package builder

import (
	"ai_agents/vision_builder/agents"
	"ai_agents/vision_builder/config"
	"ai_agents/vision_builder/model"
	"ai_agents/vision_builder/prompts"

	// "context"

	// "context"
	// "time"

	// "ai_agents/vision_builder/model"
	"ai_agents/vision_builder/env"
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	// "github.com/tmc/langchaingo/schema"
)

type VisionBuilder struct {
	ID              int
	Resume bool
	VisionStatement           string
	ActionAgentRolloutNumIter int
	ActionAgentTaskMaxRetries int
	ActionAgent               *agents.ActionAgent
	AI agents.CloseAIModelServicer
	Messages []agents.Message
	Conversations []string
	Context string
	Output          CodeResponse
	AllBuilder      strings.Builder
	UnMappedBackLog []string
	ResearchBackLog []string
	Reader          *bufio.Reader
	Fnum            int
	FilePath        string
	Goal            string
	Task            string
	SubTask         string
	Env *env.DockerClient
	EnvWaitTicks string
	BackLog agents.CurriculumAgent
	ScrumMaster *SkillManager
	AbsolutePath string
	MaxIterations int
	LastEvents string
}

func NewVisionBuilder( reader *bufio.Reader) *VisionBuilder {
	// Initialize the DockerClient system
	dc, err := env.NewDockerClient("Build any App")
	if err != nil{
		log.Fatalln(err)
	}
	//Parse the vision statement and extract tasks/sub-goals
	//Define AI configuration
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
	ckptDir := filepath.Join(currentDir, "ckpt")
	AIAgent := agents.NewActionAgent(openAI, ckptDir, false, false, false)
	AIAgent.AbsolutePath = currentDir	
	resume := false
	modelName := "gpt-3.5-turbo"
	temperature := 0
	retrievalTop_k := 5
	requestTimeout := 240
	maxIterations	:= 5
	actionAgentTaskMaxRetries := 3

	sm := NewSkillManager(modelName, float64(temperature),int32(retrievalTop_k),requestTimeout,ckptDir,resume)

	return &VisionBuilder{
		MaxIterations: maxIterations,
		AbsolutePath: currentDir,
		ActionAgent:     AIAgent,
		AI: closeAI,
		Reader: reader,
		Env: dc,
		ScrumMaster: sm,
		Resume: resume,
		ActionAgentTaskMaxRetries: actionAgentTaskMaxRetries,
	}
}

// Learn method translated from Python to Go.

func (pj *VisionBuilder) Learn(vision *model.Vision, resetEnv bool) map[string]interface{} {
	Data := struct{
		Vision string
		Program string
		Next string
		ExampleCode string
	}{
		Vision: strings.Split(vision.Description, "Enhanced Vision: ")[0],
		Program: "",
		Next: `- Provide Users and Objects dereived from the vision statement, you can use the above listed users and/or objects but do not duplicate them.
- Model Structs: generate Go code representing these users and objects as Go structs.`,
		ExampleCode: `
		...

		// PetOwner representing ... from the vision statement (For EXAMPLE)
		type PetOwner struct {
			ID        int
			Username  string
			Email     string
			Phone     string
			Pets      []Pet
			Profile   Profile
		}

		...

		// Pet ...
		type Pet struct {
			ID       int
			Name     string
			Breed    string
			Age      int
			Owner    PetOwner
		}

		...`,
	}
	humanMessage := "Generate new Models as code do not repeat models for the functionality of the app and provide what's next towards actualizing the vision"
	for{
		prompt, err := prompts.TemplateExecute(pj.AbsolutePath, "models", Data)
		if err != nil{
			log.Fatalln(err)
		}
		
		// pj.ActionAgent.
		fmt.Printf("\nTask: %s\nEnter another Task or Press Enter to continue with Task: ", humanMessage)	
		inputCode, err := pj.Reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		inputCode = strings.TrimSpace(inputCode)
		if inputCode != "" && inputCode != "Yes" && inputCode != "yes" && inputCode != "Y" && inputCode != "y"{
			humanMessage = inputCode
			
		}
		fmt.Printf("\n:::::::::input:::::::::::\n%s\n", prompt)
		output, err := pj.AI.PromptAI(prompt, humanMessage)
		fmt.Printf("\n:::::::::output:::::::::::\n%s\n", output)
		if err != nil{
			log.Fatalln(err)
		}
		blocks, section, err := ExtractCodeBlocks(output)
		if err != nil{
			log.Fatalln(err)
		} 
		if pj.ActionAgentRolloutNumIter == 0{
			pj.ActionAgent.AppMemory["models"] = MineWorkersAndModels(blocks[0])	
			fmt.Printf("\nTask: See next task above, but you can repeat previous task following below instruction:\n\nType \"Yes\", \"yes\", \"Y\" or \"y\" and Enter to Repeat last task: ")	
			inputCode, err := pj.Reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			inputCode = strings.TrimSpace(inputCode)
			if inputCode == "Yes" || inputCode == "yes" || inputCode == "Y" || inputCode == "y"{
				if aa, ok := pj.ActionAgent.AppMemory["models"].(map[string]interface{}); ok{
					for _,v := range aa{
						fmt.Println("k=",v.(map[string]model.App))
						for _, parentApp :=  range v.(map[string]model.App){
							fmt.Printf("\nParent %s", parentApp.Code)
							Data.Program += fmt.Sprintln(parentApp.Code)
							for _, child := range parentApp.Children{
								fmt.Printf("\nChild %s", child.Code)
							}
						}
					}
				}
				continue				
			}		
			Data.Next = fmt.Sprintf("- %s", strings.Split(section, "Next:")[1])
			if strings.Contains(Data.Next, "worker") || strings.Contains(Data.Next, "Worker"){
				goals := "\nGoals:"
				goalsBody := ""
				i := 1
				for _,v := range vision.Goals{
					goalsBody += fmt.Sprintf("%d. %s\n",i, v.Description)
					i++
				}
				Data.Vision = fmt.Sprintf("%s\n%s\n%s\n", Data.Vision, goals, goalsBody)
				goalsBody = ""
				pj.ActionAgentRolloutNumIter = 1
				Data.ExampleCode = `// ReviewWorker handles feedback and review systems (For EXAMPLE)
				type ReviewWorker struct {
					// any dependencies can be added here
				}
				Ok
				// LeaveReview method allows pet owners to provide feedback on services received
				func (rw *ReviewWorker) LeaveReview(serviceID int, reviewDetails map[string]interface{}) error {
					// implementation
				}`				
				humanMessage = "Define the Workers and their methods based on the goals derived from the vision listed above. Each Worker should encapsulate a set of related functionalities."
			}
		}else if pj.ActionAgentRolloutNumIter == 1{
			pj.ActionAgent.AppMemory["workers"] = MineWorkersAndModels(blocks[0])
			if (!strings.Contains(Data.Next, "worker")) && (!strings.Contains(Data.Next, "Worker")){
				for _, workersOrModels := range pj.ActionAgent.AppMemory{
					for _, workersCode := range workersOrModels.(map[string]interface{}){
						for _, copyOfParent := range workersCode.(map[string]model.App){
							fmt.Printf("\n:::::::::::App %s Start:::::::::::\n",copyOfParent.Name)
							for _, copyOfChild := range copyOfParent.Children{
								fmt.Println(":::::::::::Children Start:::::::::::::::")
								fmt.Printf("\n%v\n",copyOfChild.Code)
								fmt.Println(":::::::::::Children End:::::::::::::::")
							}
						}
					}
				}	
				// pj.ScrumMaster.Vectordb.AddDocuments(context.Background(), []schema.Document{
				// 	{PageContent: "", Metadata:  map[string]any{} },			
				// })
			}
		}
		Data.Program = blocks[0]		
	}
    return map[string]interface{}{}
}

func MineWorkersAndModels(code string) map[string]interface{} {
	// Split the code into lines.
	lines := strings.Split(code, "\n")
	fmt.Println("Code Splitted into: ", len(lines))
	var workers = make(map[string]model.App)
	var currentBlock []string
	structName := ""
	des := ""
	currentType := ""
	methodName := ""
	for cl, line := range lines {
		fmt.Println("Processing line:", cl, line)
		if line == "" {
			// When encountering an empty line, determine if we're at the end of a block.
			if len(currentBlock) > 0 {
				// Join the lines in the current block into a single string.
				block := strings.Join(currentBlock, "\n")
				currentBlock = nil // Reset the current block.
				if currentType == "type" {
					workers[structName] = model.App{
						Name:        structName,
						Description: des,
						Code:        block,
						Children:    make(map[string]model.App),
					}
				}
				if currentType == "func" {
					copyOfApp, exist := workers[structName]
					if exist {
						copyOfApp.Children[methodName] = model.App{
							Name:        methodName,
							Description: des,
							Code:        block,
						}
					}
					workers[structName] = copyOfApp
				}
				methodName = ""
			}
		} else {
			// Add the current line to the current block.
			currentBlock = append(currentBlock, line)
			if strings.HasPrefix(line, "type ") {
				structName = strings.TrimSpace(strings.Split(line, " ")[1])
				des = fmt.Sprintf("// %s ...", structName)
				if (cl > 0) && strings.Contains(lines[cl-1], "//") {
					des = strings.Split(lines[cl-1], "// ")[1]
				}
				currentType = "type"
			}
			if strings.HasPrefix(line, "func (") {
				funcPhrase := strings.TrimSpace(strings.Split(line, " ")[3])
				methodName = strings.Split(funcPhrase, "(")[0]
				des = fmt.Sprintf("// %s ...", methodName)
				if (cl > 0) && strings.Contains(lines[cl-1], "//") {
					des = strings.Split(lines[cl-1], "// ")[1]
				}
				currentType = "func"
			}
		}
	}

	// Add the last block if it's not empty.
	if len(currentBlock) > 0 {
		block := strings.Join(currentBlock, "\n")
		if currentType == "type" {
			workers[structName] = model.App{
				Name:        structName,
				Description: des,
				Code:        block,
				Children:    make(map[string]model.App),
			}
		}
		if currentType == "func" {
			copyOfApp, exist := workers[structName]
			if exist {
				copyOfApp.Children[methodName] = model.App{
					Name:        methodName,
					Description: des,
					Code:        block,
				}
			}
			workers[structName] = copyOfApp
		}
	}
	fmt.Println("Final map of workers:", workers)
	result := make(map[string]interface{})
	result["workers"] = workers
	return result
}




	
