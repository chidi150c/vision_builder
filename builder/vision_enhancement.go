package builder

import (
	"ai_agents/vision_builder/model"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	// 	"ai_agents/agile_meth/model"
	// 	"bufio"
	// 	"errors"
	// 	"fmt"
	// 	// "io/ioutil"
	// 	"log"
	// 	"os"
	// 	// "os/exec"
	// 	"regexp"
	// 	"strconv"
	// 	"strings"
	// 	"time"
)

type CodeResponse struct {
	Resoning string
	Code     string
	PreCode string
	ExecOutput string
	Report   []string
}

// func NewVisionBuilder(ai ai_model.AIModelServicer, vision *model.Vision, reader *bufio.Reader) *VisionBuilder {
// 	return &VisionBuilder{
// 		// Backlog:         NewBacklog(),
// 		AI: ai,
// 		// Vision: vision,
// 		Reader: reader,
// 	}
// }

func (pj *VisionBuilder) InProgress(ch chan bool) {
	fmt.Printf("\n")
	for {
		select {
		case <-ch:
			fmt.Printf("\n")
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c in progress...", r)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (pj *VisionBuilder) VisionEnhancement(vision *model.Vision) {
	i := 1
loop:
	for {
		progressCh := make(chan bool)
		go pj.InProgress(progressCh)
		_ = pj.ArticulatedVision(vision)
		//uncheck
		// fmt.Printf("\nIteration %d:\n\n", i)
		if i == 1 {
			_, vision = pj.BreakDownVisionIntoGoals(vision)
		} else {
			_, vision = pj.BreakDownVisionIntoNextGoals(vision)
		}
		if i == 1 {
			_, err := pj.CreateUserStoriesForTheVision(vision)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			_ = pj.ArticulatedVision(vision)
			_, err := pj.CreateNextUserStoriesForTheVision(vision)
			if err != nil {
				log.Fatalln(err)
			}
		}

		pj.MapUserStoriesToGoals(vision)
		// fmt.Printf("\nLenth of Backlog %d and lenght of goals %d", len(pj.UnMappedBackLog), len(vision.Goals))
		features := ""
		if len(pj.UnMappedBackLog) != 0 {
			for _, v := range pj.UnMappedBackLog {
				features = features + "; " + v
			}
			vision.Description, _ = pj.AI.PromptAI("Summarize in a brief paragraph of an acticulated vision satement", fmt.Sprintf("%s\nFeature: %s", vision.UpdatedVision, features))
			pj.UnMappedBackLog = []string{}
		}
		input := fmt.Sprintf("Vision: %s", vision.Description)
		for k, _ := range vision.Goals {
			input = fmt.Sprintf("%s\nFeature: %s", input, k)
		}
		vision.Description, _ = pj.AI.PromptAI("Summarize in a very brief paragraph of an acticulated vision satement", input)
		progressCh <- true
		fmt.Printf("\n\nEnhanced Vision: %s\n", vision.Description)
		for {
			fmt.Print("\nType \"Ok\" to continue or \"Redo\" to argument and press ENTER: ")
			inputCode, _ := pj.Reader.ReadString('\n')
			inputCode = strings.TrimSpace(inputCode)
			if inputCode != "Ok" && inputCode != "Redo" {
				fmt.Println("\nInvalid Input!!!!")
				continue
			} else if inputCode == "Ok" {
				break loop
			} else {
				break
			}
		}
		i++
	}
}

// func (pj *VisionBuilder) ManualRun(vision *model.Vision, reader *bufio.Reader) {
// 	r := 1
// 	t := 1
// 	y := 1

// 	chp := make(chan bool)
// 	pj.Output.Code = `
// 	package main

// 	// Worker interface defines the methods that all workers should implement
// 	type Worker interface {
		
// 	}
	
// 	func main() {

// 	}
// 	`              
// 	pj.Vision = fmt.Sprintf("Vision: %s\n", vision.Description)   
// 	for {
// 		fmt.Printf("\nGoals: \n")
// 		lengoal := len(vision.Goals)
// 		var concepts = make([]string, lengoal+1)
// 		r = 1
// 		for k, goal := range vision.Goals {
// 			concepts[r] = k
// 			fmt.Printf("%d. %s\n", r, goal.Description)
// 			fmt.Printf("   Goal Reason: %s\n", goal.GoalReasoning)
// 			fmt.Println()
// 			r++
// 		}
// 		fmt.Print("\nEnter a Goal number to implement: ")
// 		choice, _ := reader.ReadString('\n')
// 		go pj.InProgress(chp)
// 		choice = strings.TrimSpace(choice)
// 		num, err := strconv.Atoi(choice)
// 		if err != nil || num < 1 || num > lengoal {
// 			chp <- true
// 			if num == 0 {
// 				break
// 			}
// 			fmt.Println("\nInvalid Input!!!!")
// 			continue
// 		} else {
// 			chp <- true
// 			goal := vision.Goals[concepts[num]]
// 			fmt.Printf("Goal: %s\n", goal.Concept)
// 			pj.Goal = fmt.Sprintf("Goal %d: %s\n", num, goal.Description)
// 			go pj.InProgress(chp)
// 			// Derive tasks from each goal
// 			Tasks := pj.DeriveTasksFromGoal(goal, vision.Description)
// 			chp <- true
// 			for {
// 				fmt.Printf("\nTasks: \n")
// 				lengoal = len(Tasks)
// 				t = 1
// 				for _, task := range Tasks {
// 					fmt.Printf("%d. %s\n", t, task.Description)
// 					t++
// 				}
// 				fmt.Print("\nEnter a Task number to implement: ")
// 				choice, _ = reader.ReadString('\n')
// 				choice = strings.TrimSpace(choice)
// 				num, err := strconv.Atoi(choice)
// 				if err != nil || num < 1 || num > lengoal {
// 					if num == 0 {
// 						break
// 					}
// 					fmt.Println("\n Invalid Input!!!!")
// 					continue
// 				} else {
// 					task := Tasks[num-1]
// 					fmt.Printf("\nTasks: %s\n", task.Description)
// 					pj.Task = fmt.Sprintf("Tasks %d: %s\n", num, task.Description)
// 					for {       
// 						fmt.Printf("\nSub-Tasks: \n")
// 						lengoal = len(task.SubTask)
// 						y = 1
// 						for _, subtask := range task.SubTask {
// 							fmt.Printf("%d. %s\n", y, subtask)
// 							y++
// 						}
// 						fmt.Print("\nEnter a Sub-Task number to implement: ")
// 						choice, _ = reader.ReadString('\n')
// 						choice = strings.TrimSpace(choice)
// 						num, err := strconv.Atoi(choice)
// 						if err != nil || num < 1 || num > lengoal {
// 							if num == 0 {
// 								break
// 							}
// 							fmt.Println("\n  Invalid Input!!!!")
// 							continue
// 						} else {
// 							sub := task.SubTask[num-1]
// 							fmt.Printf("\nSub-Tasks: %s\n", sub)
// 							if strings.TrimSpace(sub) != "" {
// 								execOut := pj.Executor(pj.Output.Code)								
// 								pj.SubTask = fmt.Sprintf("Sub-Tasks %d: %s\n", num, sub)
// 								pj.Output.ExecOutput = fmt.Sprintf("\nExec Output: %s\n", execOut)
// 								input := fmt.Sprintf("\n%s\n%s\n%s\n%s\nCode: %s\n%s\n", pj.Vision,pj.Goal,pj.Task,pj.SubTask,pj.Output.Code,pj.Output.ExecOutput)							
// 								fmt.Printf("\n:::::::::::Input From ManualRun1:::::::::\n%s",input)								
// 								pj.Output = pj.AutoCodeDeveloper(input, task.Description)			
// 							}
// 							fmt.Print("\nEnter your next task or sub-task and Press Enter or Press Enter to continue? \n\n")
// 							inputCode, _ := reader.ReadString('\n')
// 							inputCode = strings.TrimSpace(inputCode)
// 							if inputCode != "" {
// 								execOut := pj.Executor(pj.Output.Code)
// 								pj.Output.ExecOutput = fmt.Sprintf("\nExec Output: %s\n", execOut)
// 								input := fmt.Sprintf("%s\n%s\n%s\n%s\nCode: %s\n%s", pj.Vision,pj.Goal,pj.Task,inputCode,pj.Output.Code,pj.Output.ExecOutput)					
// 								fmt.Printf("\n:::::::::::::Input From ManualRun2::::::::::\n %s",input)								
// 								pj.Output = pj.AutoCodeDeveloper(input, task.Description)
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// func (pj *VisionBuilder) CodeSummarizer(){
// 	var err error
// 	pj.Output.Code, err = pj.AI.PromptAI(CodeSummarizer2, pj.Output.Code)
// 	if err != nil{
// 		log.Fatalln(err)
// 	}
// }
// func (pj *VisionBuilder) AutomaticRun(vision *model.Vision) {
// 	pj.Output.Code = `
// 	package main

// import (
// 	"fyne.io/fyne/v2/app"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/widget"
// )

// func main() {
// 	myApp := app.New()
// 	myWindow := myApp.NewWindow("Simple Code Editor")

// 	// Create a multi-line text input
// 	codeInput := widget.NewMultiLineEntry()
// 	codeInput.SetPlaceHolder("Start typing your code...")

// 	// You can set up a menu and other features here

// 	// Add the text input to the window
// 	myWindow.SetContent(container.NewVBox(
// 		codeInput,
// 		// Add other widgets as needed
// 	))

// 	myWindow.ShowAndRun()
// }

// 	`

// 	r := 1
// 	t := 1
// 	y := 1
// 	pj.Vision = fmt.Sprintf("Vision: %s\n", vision.Description)  
// 	for {
// 		for _, goal := range vision.Goals {
// 			pj.Goal = fmt.Sprintf("Goal %d: %s\n", r, goal.Description)
// 			// input := fmt.Sprintf("\n%s\n%s\n",pj.Vision,  pj.Goal)
// 			// fmt.Printf("\n::::::::::::Input From AutomaticRun1::::::::::\n%s",input)								
// 			// pj.Output = pj.AutoDeveloper(input, goal.Description) //the first and initial composure
// 			//Derive tasks from each goal
// 			Tasks := pj.DeriveTasksFromGoal(goal, vision.Description)
// 			t = 1
// 			for _, task := range Tasks {
// 				pj.Task = fmt.Sprintf("Tasks %d: %s\n", t, task.Description)
// 				// pj.Output.ExecOutput = fmt.Sprintf("\nExec Output: %s\n", pj.Executor(pj.Output.Code))
// 				y = 1
// 				for _, sub := range task.SubTask {		
// 					pj.CodeSummarizer()			
// 					pj.SubTask = fmt.Sprintf("\n  Sub-Tasks %d: %s\n", y, sub)
// 					input := fmt.Sprintf("\n%s\n%s\n%s\n%s\nCode: %s\n%s\n", pj.Vision,pj.Goal,pj.Task, pj.SubTask, pj.Output.Code, pj.Output.ExecOutput)
// 					fmt.Printf("\n::::::::::::New SubTask Input From AutomaticRun2::::::::::\n %s",input)								
// 					pj.Output = pj.AutoCodeDeveloper(input, task.Description)					
// 					y++
// 				}
// 				t++
// 			}
// 			r++
// 		}
// 	}
// }
// func (pj *VisionBuilder) AutoDeveloper(VisionGoal, goal string) CodeResponse {
// 	chp := make(chan bool)
// 	go pj.InProgress(chp)
// 	code, _ := pj.AI.PromptAI(CodePrompt2, VisionGoal)
// 	if strings.Contains(code, "No-code") {
// 		chp <- true
// 		fmt.Printf("\nCurrent task requires your action: %s\n\n", code)
// 		return CodeResponse{}
// 	}
// 	codeBlocks, report, err := extractCodeBlocks(code)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	if len(codeBlocks) > 1 {
// 		sb := strings.Builder{}
// 		for _, v := range codeBlocks{
// 			sb.WriteString(fmt.Sprintln(v))
// 		}
// 		pj.Output.Code = fmt.Sprintf("\n%s\n", sb.String())
// 		input := fmt.Sprintf("\n%s\n%s\nCode: %s\n",pj.Vision,  pj.Goal, pj.Output.Code)
// 		fmt.Printf("\n::::::::::::::::Input From AutoDeveloper::::::::::::\n %s",input)								
// 		pj.AutoCodeDeveloper(input, goal)
// 	}
// 	pj.Output.Resoning = report
// 	// Define the file path
// 	chp <- true
// 	pj.CodeOnFile(pj.Output.Code)
// 	return pj.Output
// }

// // func (pj *VisionBuilder)Executor(code string)string {
// // 	// Step 1: Write the Go code to a temporary file
// // 	tmpFile, err := ioutil.TempFile("", "*.go")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}
// // 	defer os.Remove(tmpFile.Name()) // Clean up

// // 	if _, err := tmpFile.Write([]byte(code)); err != nil {
// // 		tmpFile.Close()
// // 		log.Fatal(err)
// // 	}
// // 	tmpFile.Close()

// // 	// Step 2: Use exec.Command to run 'go run <tempfile>'
// // 	cmd := exec.Command("go", "run", tmpFile.Name())
// // 	byteOut, err := cmd.CombinedOutput()
// // 	output := string(byteOut)
// // 	if err != nil {
// // 		output = fmt.Sprintf("%s: Failed to execute 'go run': %v \n", output,  err)
// // 	}
// // 	// Step 3: Process output
// // 	return output
// // }

// func extractCodeBlocks(output string) ([]string, string, error) {
// 	// Adjusting the regex to match different newline characters and be more robust
// 	// This pattern is more flexible in capturing the content between ```go and ```
// 	pattern := "(?s)```go\r?\n(.*?)[\r\n]```"
// 	re := regexp.MustCompile(pattern)
// 	sections := strings.Split(output, "```go")
// 	// Find all matches
// 	matches := re.FindAllStringSubmatch(output, -1)

// 	if len(matches) == 0 {
// 		fmt.Println("No code blocks found. Ensure the file contains correctly delimited Go code blocks.")
// 		return nil, "", errors.New("no code blocks found")
// 	}

// 	// Print each code block found
// 	blocks := make([]string, len(matches))
// 	// var sb strings.Builder
// 	for i, blockOfLines := range matches {
// 		blocks[i] = blockOfLines[1]
// 	}
// 	fmt.Printf("\n\n::::::::Reoprt From Executor:::::::::\nTotal number of Code blocks generated: \n%d\n\n%s\n", len(blocks), sections[0])
// 	return blocks, sections[0], nil
// }
// func (pj *VisionBuilder) AutoCodeDeveloper(visionGoalAndTask, task string) CodeResponse {
// 	pj.Output.PreCode = pj.Output.Code
// 	// i := 1
// 	for {
// 		chp := make(chan bool)
// 		go pj.InProgress(chp)
// 		visionGoalAndTask, _ = pj.AI.PromptAI(CodePrompt3, visionGoalAndTask)
// 		if strings.Contains(visionGoalAndTask, "No-code") {
// 			chp <- true
// 			fmt.Printf("\nCurrent task requires your action: %s\n\n", visionGoalAndTask)
// 			return CodeResponse{}
// 		}
// 		//Extract the code sections from the AI model output
// 		codes, _, err := extractCodeBlocks(visionGoalAndTask)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
	
// 		chp <- true
// 		if len(codes) >= 1 {
// 			sb := strings.Builder{}
// 			for _, v := range codes {
// 				sb.WriteString(fmt.Sprintf("\n%s\n", v))
// 			}
// 			pj.Output.Code = fmt.Sprintf("\n%s\n", sb.String())	
// 		}else{
// 			panic(visionGoalAndTask)
// 		}
// 		execOut := pj.Executor(pj.Output.Code)			
// 		pj.Output.ExecOutput = fmt.Sprintf("\nExec Output: %s\n", execOut)
// 		visionGoalAndTask = fmt.Sprintf("\n%s\n%s\n%s\n%s\nCode: %s\n%s\n",pj.Vision,  pj.Goal, pj.Task, pj.SubTask, pj.Output.Code, pj.Output.ExecOutput)
// 		if strings.Contains(execOut, "Failed") {
// 			fmt.Printf("::::::::::::Looping Back Inside AutoCodeDeveloper::::::::::::\n%s", visionGoalAndTask)
// 			if pj.Fnum > 0{
// 				pj.Fnum--
// 			}				
// 			pj.CodeOnFile(pj.Output.Code)			
// 			fmt.Println("\nCode Failled Interven and/or Press Enter to Continue: ")
// 			human, _ := pj.Reader.ReadString('\n')
// 			pj.Output.ExecOutput = pj.Output.ExecOutput+" \n My suggestion to address it: "+human
// 			visionGoalAndTask = fmt.Sprintf("\n%s\n%s\n%s\n%s\nCode: %s\n%s\n",pj.Vision,  pj.Goal, pj.Task, pj.SubTask, pj.Output.Code, pj.Output.ExecOutput)
// 			// i++
// 			continue
// 		} else {
// 			fmt.Printf("::::::::::::Moving On To Next Task Inside AutoCodeDeveloper::::::::::::\n%s", visionGoalAndTask)
// 			break
// 		}
// 	}
// 	pj.CodeOnFile(pj.Output.Code)
// 	fmt.Println("\nPress Enter to continue: ")
// 	_, _ = pj.Reader.ReadString('\n')
// 	return pj.Output
// }
// func (pj *VisionBuilder)CodeOnFile(code string) {
// 	pj.FilePath = fmt.Sprintf("../output/app%d.go", pj.Fnum)
// 	pj.Fnum++
// 	// Use os.OpenFile with the appropriate flags
// 	// os.O_CREATE - Create the file if it does not exist
// 	// os.O_WRONLY - Open the file] write-only
// 	// os.O_TRUNC - If possible, truncate the file when opened
// 	file, err := os.OpenFile(pj.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
// 	if err != nil {
// 		log.Fatalf("Failed to open or create file: %s", err)
// 	}
// 	// Write the content to the file
// 	if _, err := file.WriteString(code); err != nil {
// 		log.Fatalf("Failed to write to file: %s", err)
// 	}
// 	file.Close()
// 	fmt.Printf("\nCode written successfully to file \"%s\".\n", pj.FilePath)
// }

// ArticulatedVision where Questions were raised to articulate the vision and answers were provided that were converted later
// as features added to the vision in order to prodcuce an Enhanced vision statement. No vision praramter was modofied.
func (pj *VisionBuilder) ArticulatedVision(vision *model.Vision) string {
	pj.AllBuilder.Reset()
	goslsForQuestion := pj.QuestionsForVisionClarification(vision)
	vi := ""
	i := 1
	pj.AllBuilder.WriteString("\n\nBELOW ARE QUESTIONS AND ANSWERS TO CLARIFY THE FEATURES AND CHALLENGES EXPECTED FOR THE VISION")
	for k, v := range goslsForQuestion {
		pj.AllBuilder.WriteString(fmt.Sprintf("\n\nQuestion %d: %s\nGoal %d: %s \n", k+1, v.Question, k+1, v.Description))
		v.Answer = pj.AnswerToQuestionForVisionClarity(v.Question, vision)
		pj.AllBuilder.WriteString(fmt.Sprintf("Answer %d: %s", k+1, v.Answer))
		vi = fmt.Sprintf("Feature %d: %s", i, v.Answer)
		if strings.Contains(vi, "Unknown:") {
			//send vi to R&D
			pj.ResearchBackLog = append(pj.ResearchBackLog, fmt.Sprintf("%s\n%s", v.Question, vi))
		} else {
			vision.Description = fmt.Sprintf("%s\n%s", vision.Description, vi)
			i++
		}
	}
	pj.AllBuilder.WriteString("\n\nTHE FOLLOWING IS THE UPDATED VISION WITH SOME FEATURES:\n")
	pj.AllBuilder.WriteString(fmt.Sprintf("\n%s\n\n", vision.Description))
	return vision.Description
}

// CreateUserStory creates a new user story with the given description and priority.
func (pj *VisionBuilder) CreateUserStoriesForTheVision(vision *model.Vision) ([]model.UserStory, error) {
	PreSalesPrompt := `
	TASK DESCRIPTION:
	You are a helpful Agile Development Team that helps to Identify User Activities and Creating User Stories for a given vision and Prioritizing the user stories based on their importance and impact on achieving the VisionBuilder goals, considering factors such as user needs, feasibility, and alignment with the VisionBuilder vision.
	
	CRITERIA:
	1. Identify all the activities that users will perform while using the Product of the VISION and for each user activity, break them down into smaller actionable items called user stories each preceded with "As a user," as label.
	2. We then prioritize these user stories based on their importance and impact on achieving the VisionBuilder goals, considering factors such as user needs, feasibility, and alignment with the VisionBuilder vision: with "Priority:" as label.
	
	RESPONSE FORMAT:
	For each user story:
	- The user story
	- Its Priority

	EXAMPLE OUTPUT: 
	1. As a user, I want to experience immersive graphics in the game so that I can be fully engaged in the virtual environment."
	Priority: High

	2. As a user, I want the game to have accessibility features for players with disabilities so that everyone can enjoy the game.
	Priority: Medium

	...

	8. As a user, I want diverse cultural elements and representation in the game world, to feel connected and represented.    
	Priority: Low

	...`
	var (
		err         error
		description string
	)
	description, err = pj.AI.PromptAI(PreSalesPrompt, vision.Description)
	if err != nil {
		return nil, err
	}
	userStories := pj.fillUserStory(description, vision)
	pj.AllBuilder.WriteString(fmt.Sprintf("THE FOLLOWING ARE %d USER STORIES OF THE VISION: %s\n\n", len(userStories), vision.UpdatedVision))
	pj.AllBuilder.WriteString("\nUser Stories:\n")
	for k, v := range userStories {
		pj.AllBuilder.WriteString(fmt.Sprintf("%d. As a user, %s\n", k+1, v.Description))
	}
	pj.AllBuilder.WriteString("\n")
	vision.DraftedUserStories = userStories
	// fmt.Println((pj.AllBuilder.String()))
	return vision.DraftedUserStories, nil
}

// CreateUserStory creates a new user story with the given description and priority.
func (pj *VisionBuilder) CreateNextUserStoriesForTheVision(vision *model.Vision) ([]model.UserStory, error) {
	PreSalesPrompt := `
	TASK DESCRIPTION:
	You are a helpful Agile Development Team that helps to Identify User Activities and Creating User Stories for a given vision and Prioritizing the user stories based on their importance and impact on achieving the VisionBuilder goals, considering factors such as user needs, feasibility, and alignment with the VisionBuilder vision.

	I will provide list of already considered user story concepts, so that you generate the next user story whose concept is not yet considered and does not align with any listed concept.
	
	INPUT FORMAT:
	Vision: [I will place the vision statement here]
	Concept [already used number]: [Already considered concept] 

	CRITERIA:
	
	1. Considering the listed concepts, identify user activities under concepts that do not align with any listed concept. 
	2. Then ensure identified user activities are based on user needs, feasibility, and alignment with the VisionBuilder vision.
	3. And for each user activity, break them down into smaller actionable items called user stories each preceded with "As a user," as label. 
	3. Then prioritize these user stories based on their importance and impact on achieving the VisionBuilder goals, considering factors such as user needs, feasibility, and alignment with the VisionBuilder vision: with "Priority:" as label.
	
	RESPONSE FORMAT:
	For each user story:
	- The user story
	- Its Priority

	EXAMPLE OUTPUT:
	1. As a user, I want to experience immersive graphics in the game so that I can be fully engaged in the virtual environment."
	Priority: High

	2. As a user, I want the game to have accessibility features for players with disabilities so that everyone can enjoy the game.
	Priority: Medium

	...

	8. As a user, I want diverse cultural elements and representation in the game world, to feel connected and represented.    
	Priority: Low

	...`
	var (
		err         error
		description string
	)
	description, err = pj.AI.PromptAI(PreSalesPrompt, vision.Description)
	if err != nil {
		return nil, err
	}
	userStories := pj.fillUserStory(description, vision)
	fmt.Printf("\nnew user experiences determined...\n")
	for k, v := range userStories {
		_ = k
		vision.DraftedUserStories = append(vision.DraftedUserStories, v)
	}
	return vision.DraftedUserStories, nil
}

func (pj *VisionBuilder) MapUserStoriesToGoals(vision *model.Vision) {
	pj.AllBuilder.Reset()
	pj.AllBuilder.WriteString("Goals:\n")
	for _, v := range vision.DraftedGoals {
		pj.AllBuilder.WriteString(fmt.Sprintf("%d. %s\n", v.ID, v.Description))
	}
	pj.AllBuilder.WriteString("\n")
	pj.AllBuilder.WriteString("\nUser Stories:\n")
	for _, v := range vision.DraftedUserStories {
		pj.AllBuilder.WriteString(fmt.Sprintf("%d. As a user, %s\n", v.ID, v.Description))
	}
	pj.AllBuilder.WriteString("\n")

	input := pj.AllBuilder.String()

	MappingPrompt := `
	TASK DESCRIPTION:
	Generate mappings between given user stories and goals outlined for a software development VisionBuilder, including reasoning for each mapping. Explicitly include the reasoning with a 'Reasoning:' label. Additionally, identify and label the underlying concepts that guide the VisionBuilder's development efforts with a 'Concept:' label.

	INPUT FORMAT:
	Provide a list of user stories and goals. Each user story and goal should be accompanied by a brief description.

	CRITERIA:
	The model should map each user story to the corresponding goal from the provided list and explicitly provide the reasoning behind the mapping with a 'Reasoning:' label. If a user story does not directly align with any goal, the model should indicate so with an explanation also preceded by the 'Reasoning:' label. Additionally, the model should identify and label the underlying concepts that guide the VisionBuilder's development efforts with a 'Concept:' label.

	RESPONSE FORMAT:
	For each mapping, the model's response should include:
	- The user story and the corresponding goal.
	- A 'Reasoning:' label followed by the rationale for the mapping.
	- A 'Concept:' label followed by the underlying concept that guides the VisionBuilder's development efforts.
	- If a user story does not align with any goal, it should be flagged as an isolated user story with reasonings for the lack of alignment, also following the 'Reasoning:' label.
	- if a "User Story 2" maps to "Goal 3" and "Goal 6" indicate only "Goal 3", that is, "User Story 2" maps to "Goal 3"
	- if "User Story 9" and "User Story 10" map to "Goal 6", separate each "User Story [x]" map to "Goal 6" into a separate line

	EXAMPLE INPUT:
	User Stories:
	1. As a user, I want to be able to select between trend following and mean reversion strategies so that I can optimize my cryptocurrency trading based on my preferred investment style and market conditions.
	2. As a user, I want to set the time horizon for my trades (short-term, medium-term, or long-term) to align with my risk tolerance and profit goals.
	Goals:
	1. Develop a trading system that supports both trend following and mean reversion strategies.
	2. Integrate trading options for various timeframes, including short-term, medium-term, and long-term strategies.

	EXAMPLE OUTPUT:
	1. User Story 1 maps to Goal 1.
	Reasoning: The desire to select between trading strategies directly supports the goal of developing a trading system with versatile strategy options.
	Concept: Versatility in trading strategies.

	...
	
	3. User Story 2 does not directly align with any specific goal.
	Reasoning: The request for personalized recommendations suggests a need for an adaptive AI component, which is not explicitly covered in the outlined goals. This might indicate a new area for development or an extension of existing goals.
	Concept: Adaptive AI recommendations.

	ADDITIONAL INSTRUCTIONS:
	Ensure the reasoning is clear, concise, and directly connects the user's needs or desires with the goals' intent and outcomes. Use the 'Reasoning:' label to clearly separate the rationale from the mapping for easy reading and comprehension. Similarly, use the 'Concept:' label to explicitly identify the underlying concepts guiding the VisionBuilder's development efforts.
`
	output, _ := pj.AI.PromptAI(MappingPrompt, input)
	// Split the output into individual mappings
	mappings := strings.Split(output, "\n\n")
	// Initialize a slice to store parsed mappings
	var storyMappings []model.UserStory
	var goalMappings []model.Goal
	// Parse each mapping
	for _, mapping := range mappings {
		lines := strings.Split(mapping, "\n")
		if len(lines) < 3 {
			fmt.Printf("\n\n\nStrange Mapping Omitted!!! %s\n\n\n", mapping)
			continue
		}
		// Extract user story, goal, reasoning, and concept from each mapping
		reason := strings.Split(lines[1], "Reasoning:")
		concepM := strings.Split(lines[2], "Concept:")
		usrgoal := strings.Split(lines[0], ".")
		if len(reason) <= 1 || len(concepM) <= 1 || len(usrgoal) <= 1 {
			continue
		}
		reasoning := strings.TrimSpace(reason[1])
		concept := strings.TrimSpace(concepM[1])
		ug := usrgoal[1]
		if strings.Contains(ug, "does not directly align") {
			pj.UnMappedBackLog = append(pj.UnMappedBackLog, fmt.Sprintf("%s feature\n", concept))
		} else {
			reGoal, _ := regexp.Compile(`Goal (\d+)`)
			reUStory, _ := regexp.Compile(`User Story (\d+)`)
			gid := reGoal.FindStringSubmatch(ug)
			uid := reUStory.FindStringSubmatch(ug)
			if len(gid) > 1 && len(uid) > 1 {
				nuid, err := strconv.Atoi(uid[1])
				if err != nil {
					log.Fatalln(err)
				}
				ngid, err := strconv.Atoi(gid[1])
				if err != nil {
					log.Fatalln(err)
				}
				// Populate the Mapping struct
				storyMappings = append(storyMappings, model.UserStory{
					ID:           nuid,
					MappedGoals:  ngid,
					MapReasoning: reasoning,
					Concept:      concept,
				})
				// Populate the Mapping struct
				goalMappings = append(goalMappings, model.Goal{
					ID:                ngid,
					MappedUserStories: nuid,
					MapReasoning:      reasoning,
					Concept:           concept,
				})
			}
		}

	}

	// Print parsed mappings
	for _, m := range storyMappings {
		for k, v := range vision.DraftedUserStories {
			if _, ok := vision.UserStories[m.Concept]; (!ok) && m.ID == v.ID {
				vision.DraftedUserStories[k].MapReasoning = m.MapReasoning
				vision.DraftedUserStories[k].Concept = m.Concept
				vision.DraftedUserStories[k].MappedGoals = m.MappedGoals
				vision.UserStories[m.Concept] = &vision.DraftedUserStories[k]
			}
		}
		// fmt.Printf("\nUser Story: %d\nGoal: %d\nReasoning: %s\nConcept: %s\n\n", m.ID, m.MappedGoals, m.Reasoning, m.Concept)
		// pj.AllBuilder.WriteString(fmt.Sprintf("User Story: %d\nGoal: %s\nReasoning: %s\nConcept: %s\n\n", m.ID, m.MappedGoals, m.Reasoning, m.Concept))
	}
	for _, m := range goalMappings {
		for k, v := range vision.DraftedGoals {
			if _, ok := vision.Goals[m.Concept]; (!ok) && m.ID == v.ID {
				vision.DraftedGoals[k].MapReasoning = m.MapReasoning
				vision.DraftedGoals[k].Concept = m.Concept
				vision.DraftedGoals[k].MappedUserStories = m.MappedUserStories
				vision.Goals[m.Concept] = &vision.DraftedGoals[k]
			}
		}
		// fmt.Printf("\nGoal: %d\nUser Story: %d\nReasoning: %s\nConcept: %s\n\n", m.ID, m.MappedUserStories, m.Reasoning, m.Concept)
		// pj.AllBuilder.WriteString(fmt.Sprintf("User Story: %d\nGoal: %s\nReasoning: %s\nConcept: %s\n\n", m.ID, m.MappedGoals, m.Reasoning, m.Concept))
	}
}

// DeriveTasksFromGoal simulates the process of deriving tasks from the given goal.
func (pj *VisionBuilder) DeriveTasksFromGoal(goal *model.Goal, visionStatement string) []*model.Task {
	// Prompt AI model to derive tasks based on the goal

	TasksPrompt := `TASK DESCRIPTION:
	Given: ` + visionStatement + `

	Break down the provided Goal to achieve the vision, into smaller actionable tasks or sub-goals, that can be easily managed and implemented by an Agile Development Team. This will help in planning, tracking progress, and ensuring that each component of the goal is addressed effectively.

	INPUT FORMAT:
	[I, not you, will provide the Goal to be broken down into tasks]

	CRITERIA:
	- You will break down the goal into smaller, actionable tasks or sub-goals.
	- Tasks or sub-goals must be clear, concise, and directly related to the original goal.
	- Include any relevant details or considerations needed to understand and implement the task.

	RESPONSE FORMAT:
	- List the smaller, actionable tasks or sub-goals.
	- Provide a brief description of each task or sub-story.
	
	EXAMPLE OUTPUT:
	1. Design a login interface.
	2. Implement authentication mechanism (e.g., OAuth, JWT).
	3. Develop multi-factor authentication feature.
	4. Create user session management.
	5. Implement encryption for user credentials.
	6. Test the login process for security vulnerabilities.
	...

	ADDITIONAL INSTRUCTIONS:
	- Ensure that the breakdown of the goal is comprehensive and covers all necessary aspects to achieve the goal's objective.
	- Consider the technical and resource constraints that may influence the implementation of tasks.
	- Use clear and understandable language to ensure that all team members can easily grasp the tasks and their purposes.
	`

	response, err := pj.AI.PromptAI(TasksPrompt, goal.Description)
	if err != nil {
		log.Fatalln(err)
	}
	// Split the response into individual tasks
	taskDescriptions := strings.Split(response, "\n\n")
	var tasks []*model.Task
	for _, desc := range taskDescriptions {
		subtask := strings.Split(strings.TrimSpace(desc), ":")
		if len(subtask) <= 1 {
			continue
		}

		task := subtask[0]
		if !strings.Contains(task, ".") {
			continue
		}
		task = strings.Split(task, ". ")[1]

		subLines := strings.Split(subtask[1], "\n")
		var newsubtks []string
		for _, subtask := range subLines {
			if !strings.Contains(subtask, "-") {
				continue
			}
			subtask = strings.Split(subtask, "-")[1]
			newsubtks = append(newsubtks, subtask)
		}

		tasks = append(tasks, &model.Task{
			Description: task,
			SubTask:     newsubtks,
		})
	}
	goal.Tasks = tasks
	return tasks
}

func (pj *VisionBuilder) BreakDownUserStoryIntoTasks(vision *model.Vision) {
	var TaskBuilder strings.Builder
	TaskBuilder.WriteString(fmt.Sprintf("Vision: %s\n\n", vision.UpdatedVision))
	k := 1
	for _, v := range vision.UserStories {
		TaskBuilder.WriteString(fmt.Sprintf("User Story %s: %s\n\n", k, v.Description))
	}

	TasksPrompt := `TASK DESCRIPTION:
	Break down provided user stories into smaller, actionable tasks or sub-stories that can be easily managed and implemented by an Agile Development Team. This will help in planning, tracking progress, and ensuring that each component of the user story is addressed effectively.

	INPUT FORMAT:
	Provide a list of user stories. Each user story should be clearly stated, focusing on what the end-user wants and why. User stories should be formatted and numbered for easy reference.

	CRITERIA:
	- Each user story should be broken down into smaller, actionable tasks or sub-stories.
	- Tasks or sub-stories must be clear, concise, and directly related to the original user story.
	- Include any relevant details or considerations needed to understand and implement the task.
	- Prioritize tasks or sub-stories that are critical for the user story's completion.
	- Identify dependencies between tasks or sub-stories where applicable.

	RESPONSE FORMAT:
	For each user story:
	- List the smaller, actionable tasks or sub-stories.
	- Provide a brief description of each task or sub-story.
	- If applicable, note any dependencies between tasks.
	- Highlight any tasks that are particularly critical for achieving the goal of the user story.

	EXAMPLE INPUT:
	User Story 1: As a user, I want to securely log in to the application so that I can access my personal dashboard.

	User Story 2: As an admin, I want to be able to review user activity logs to ensure system security and compliance.

	EXAMPLE OUTPUT:
	Breakdown for User Story 1:
		1. Design a login interface.
		2. Implement authentication mechanism (e.g., OAuth, JWT).
		3. Develop multi-factor authentication feature.
		4. Create user session management.
		5. Implement encryption for user credentials.
		6. Test the login process for security vulnerabilities.

	Breakdown for User Story 2:
		1. Develop functionality to capture user activity logs.
		2. Design and implement a log review interface for admins.
		3. Integrate filtering and search capabilities for logs.
		4. Implement access controls to ensure only admins can review logs.
		5. Test log review features for usability and security.

	ADDITIONAL INSTRUCTIONS:
	- Ensure that the breakdown of each user story is comprehensive and covers all necessary aspects to achieve the user story's objective.
	- Consider the technical and resource constraints that may influence the implementation of tasks.
	- Use clear and understandable language to ensure that all team members can easily grasp the tasks and their purposes.
	`
	tasks, err := pj.AI.PromptAI(TasksPrompt, "")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("TTTTTaslks=: ", tasks)
	tasksOfEachUserStories := strings.Split(tasks, "\n\n")
	for k, v := range tasksOfEachUserStories {
		_, _ = k, v
		// fmt.Println(v)
		// tasksOfAUserStory := strings.Split(v, "\n")
		// for y, v := range tasksOfAUserStory {
		// 	fmt.Println(v)
		// 	vision.UserStories[k].Tasks[]&model.Task{
		// 		ID:          y,
		// 		Description: v,
		// 	}
		// }
	}
}

func (pj *VisionBuilder) QuestionsForVisionClarification(vision *model.Vision) []*model.Goal {
	QuestionPrompt := `
    TASK DESCRIPTION:
    As my diligent assistant, your role is pivotal in driving and determining the goals necessary to actualize the vision I've provided. Your questions will serve to guide us towards the realization of the vision.

    INPUT FORMAT:
    Vision: [My vision comes here]

    CRITERIA:
	1. To clearly understand the vision or identify features and/or options specific to actualizing the vision, ask relevant questions that leads to realization of the vision.
    2. Each question should be aimed at clarifying specific aspects of the vision or propose a potential objectives that aligns with the goals of the vision..
    3. Pose a minimum of 5 and not more than 10 questions that directly contribute to identifying goals crucial for the actualization of the provided vision.
    4. Focus solely on generating insightful inquiries to guide the goal-setting process.

    RESPONSE FORMAT:
    Reasoning: [Your reasoning behind the question that contributes to identifying crucial goals for actualizing the vision]
    Question 1: [Your question here]
    Question 2: [Your question here]
    ...
    
    ADDITIONAL INSTRUCTIONS:
    Ensure your reasoning is directly related to the vision statement and highlights the importance of the questions in guiding the goal-setting process.

	`
	goals, err := pj.AI.PromptAI(QuestionPrompt, vision.Description)
	if err != nil {
		log.Fatalln(err)
	}
	gl := fillClarityQuestions(goals)
	return gl
}

// AnswerToQuestionForVisionClarity function generates an answer for the question using LLM.
func (pj *VisionBuilder) AnswerToQuestionForVisionClarity(question string, vision *model.Vision) string {
	input := fmt.Sprintf("Vision : %s\nQuestion : %s", vision.Description, question)
	AnswersPrompt := `
	TASK DESCRIPTION:
	As my dedicated assistant, I'm seeking your expertise to clarify the vision for our VisionBuilder. Your insightful answers to questions related to the VisionBuilder's vision will guide us through this process.

	CRITERIA:
	1. I'll share the VisionBuilder's vision statement with you and pose a question based on it.
	2. I trust your judgment to select the most efficient approach from multiple approaches/options to achieve the vision.
	3. Ensure your responses are very brief and begin with "Answer:" for clarity.
	4. Your answers should reflect a proactive approach towards achieving the vision.
	5. In situations where there are several approaches/choices to decision/choose from, decide on or choose a specific approach/choice that would lead us towards achieving the vision fastest with less resources.
	6. In situations where specific information is lacking, simply respond with "Answer: Unknown:".
	7. In situations where no specific decision/choice can be made, simply respond with "Answer: Unknown:" and number out the choices to select from.
	8. Ensure that every response contributes directly to realizing the VisionBuilder's vision.
	
	RESPONSE FORMAT:
	Answer: [Unknown:][Your answer here][Numbered choices]
	`
	qaAnswer, err := pj.AI.PromptAI(AnswersPrompt, input)
	if err != nil {
		log.Fatalln(err)
	}
	return strings.Split(qaAnswer, "Answer: ")[1]
}

// Function to break down pj vision into 5 goals
func (pj *VisionBuilder) BreakDownVisionIntoNextGoals(vision *model.Vision) ([]model.Goal, *model.Vision) {
	ScrumPrompt := `
	TASK DESCRIPTION:
	You are a helpful assistant that tells me the NEXT goals required to achieve a given VisionBuilder Vision. The ultimate goal is to actualize the Vision by accoplishing these goals. Clearly articulate the Vision for the VisionBuilder and outline specific goals; what to do in order to achieve the Vision. 
	
	I will provide already considered goal concepts, so that you generate the next goals whose concept is not yet considered

	INPUT FORMAT:
	Vision: [I will place the vision statement here]
	Concept [already used number]: [Already considered concept] 
	
	CRITERIA:
	1. You should act as an agile scrum master.
	2. Clearly articulate the Vision for the VisionBuilder.
	3. Outline specific goals, whose concept does not relate to any listed concepts, but are required in order to achieve the Vision.
	4. Generate atleast 5 Unique NEXT Goals and not more than 10 NEXT Goals, each in a breif sentence.
	5. Do not include a goal whose concept is already considered
 
	RESPONSE FORMAT:
	Vision: [Place the clearly articulated vision statement here]
	Reasoning 1: [Your reasoning for the next goal here]
	Goal 1: [Your next goal to realize the vision here]
	Reasoning 2: [Your reasoning for the next goal here]
	Goal 2: [Your next goal to realize the vision here]
	Reasoning 3: [Your reasoning for the next goal here]
	Goal 3: [Your next goal to realize the vision here]
	...`

	input := fmt.Sprintf("Vision: %s", vision.Description)
	for k, v := range vision.Goals {
		input = fmt.Sprintf("%s\nConcept %d: %s", input, v.ID, k)
	}
	//uncheck
	// fmt.Printf("%v", input)
	goals, err := pj.AI.PromptAI(ScrumPrompt, input)
	if err != nil {
		log.Fatalln(err)
	}
	vgoals, updateVision := fillGoalandReasoning(goals, vision)
	vision.UpdatedVision = updateVision //This is a clearly acticulated vision
	fmt.Printf("\nnew goals determined...\n")
	for k, v := range vgoals {
		_ = k
		vision.DraftedGoals = append(vision.DraftedGoals, v)
	}
	return vision.DraftedGoals, vision
}

// Function to break down pj vision into 5 goals
func (pj *VisionBuilder) BreakDownVisionIntoGoals(vision *model.Vision) ([]model.Goal, *model.Vision) {
	ScrumPrompt := `
	TASK DESCRIPTION:
	You are a helpful assistant that tells me the goals required to achieve a given VisionBuilder Vision. The ultimate goal is to actualize the Vision by accoplishing these goals. Clearly articulate the Vision for the VisionBuilder and outline specific goals; what to do in order to achieve the Vision.

	INPUT FORMAT:
	Vision: [I will place the vision statement here]
	
	CRITERIA:
	1. You should act as an agile scrum master.
	2. Clearly articulate the Vision for the VisionBuilder.
	3. Outline specific goals; what to do in order to achieve the Vision.
	4. Generate atleast 5 Unique Goals and not more than 10 Goals, each in a breif sentence.
 
	RESPONSE FORMAT:
	Vision: [Place updated vision statement here]
	Reasoning 1: [Your reasoning for the first goal here]
	Goal 1: [Your first goal to realize the vision here]
	Reasoning 2: [Your reasoning for the next goal here]
	Goal 2: [Your next goal to realize the vision here]
	Reasoning 3: [Your reasoning for the next goal here]
	Goal 3: [Your next goal to realize the vision here]
	...`

	goals, err := pj.AI.PromptAI(ScrumPrompt, vision.Description)
	if err != nil {
		log.Fatalln(err)
	}
	vgoals, updateVision := fillGoalandReasoning(goals, vision)
	vision.UpdatedVision = updateVision //This is a clearly acticulated vision
	vision.DraftedGoals = vgoals
	return vgoals, vision
}

func fillClarityQuestions(input string) []*model.Goal {
	// Define regular expressions to match question and reasoning patterns.
	questionRegex := regexp.MustCompile(`Question \d+: (.+)`)
	reasoningRegex := regexp.MustCompile(`Reasoning : (.+)`)

	// Slice to store the extracted questions and reasoning.
	var visionQuestions []*model.Goal

	// Split the input string by newlines to process each line separately.
	lines := strings.Split(input, "\n\n")

	// Iterate over each line in the input string.
	for _, line := range lines {
		// Try to find a question and reasoning in the line.
		questionMatch := questionRegex.FindStringSubmatch(line)
		reasoningMatch := reasoningRegex.FindStringSubmatch(line)

		// If both a question and reasoning are found in the line, add them to the visionQuestions slice.
		if len(questionMatch) > 1 && len(reasoningMatch) > 1 {
			visionQuestions = append(visionQuestions, &model.Goal{
				Question:          questionMatch[1],
				QuestionReasoning: reasoningMatch[1],
			})
		} else if len(questionMatch) > 1 {
			// If only a question is found, add it to the visionQuestions slice without reasoning.
			visionQuestions = append(visionQuestions, &model.Goal{
				Question: questionMatch[1],
			})
		}
	}
	return visionQuestions
}

// Function to parse the text and fill the VisionGoal struct
func fillGoalandReasoning(input string, vision *model.Vision) ([]model.Goal, string) {
	var (
		visionGoals    []model.Goal
		reasoningMatch []string
		goalMatch      []string
		visionMatch    []string
	)
	visionRegex := regexp.MustCompile(`Vision: (.+)`)
	reasoningRegex := regexp.MustCompile(`Reasoning \d+: (.+)`)
	goalRegex := regexp.MustCompile(`Goal \d+: (.+)`)
	visionMatch = visionRegex.FindStringSubmatch(input)
	lines := strings.Split(input, "\n\n")
	for _, line := range lines {
		reasoningMatch = reasoningRegex.FindStringSubmatch(line)
		goalMatch = goalRegex.FindStringSubmatch(line)
		if reasoningMatch != nil && goalMatch != nil {
			visionGoals = append(visionGoals, model.Goal{
				ID:            <-vision.NextGoalIDChan,
				GoalReasoning: strings.TrimSpace(reasoningMatch[1]),
				Description:   strings.TrimSpace(goalMatch[1]),
			})
		}
	}
	vi := vision.Description
	if len(visionMatch) > 1 {
		vi = strings.TrimSpace(visionMatch[1])
	}
	return visionGoals, vi
}

func (pj *VisionBuilder) fillUserStory(description string, vision *model.Vision) []model.UserStory {
	var (
		userStories    []model.UserStory
		userStoryMatch []string
		priorityMatch  []string
	)
	userStoryRegex := regexp.MustCompile(`As a user, (.+)`)
	priorityRegex := regexp.MustCompile(`Priority: (.+)`)
	lines := strings.Split(description, "\n\n")
	for _, line := range lines {
		userStoryMatch = userStoryRegex.FindStringSubmatch(line)
		priorityMatch = priorityRegex.FindStringSubmatch(line)
		if userStoryMatch != nil && priorityMatch != nil {
			userStories = append(userStories, model.UserStory{
				ID:          <-vision.NextUserStoryIDChan,
				Description: strings.TrimSpace(userStoryMatch[1]),
				Priority:    strings.TrimSpace(priorityMatch[1]),
			})
		}
	}
	return userStories
}
