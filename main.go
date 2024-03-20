package main

import (
	"ai_agents/vision_builder/builder"
	"ai_agents/vision_builder/model"
	"bufio"
	"log"

	// "ai_agents/vision_builder/env"
	"flag"
	"fmt"
	"os"
	"strings"
)
func validateVisionStatement(visionStatement *string) error {
    if *visionStatement == "" {
        // Prompt the user to enter the vision statement
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter the vision statement: ")
        statement, err := reader.ReadString('\n')
        if err != nil {
            return fmt.Errorf("error reading input: %v", err)
        }
        *visionStatement = strings.TrimSpace(statement)
    }

    // Validate input
    if *visionStatement == "" {
        return fmt.Errorf("vision statement is required")
    }

    return nil
}
func main() {
	vis := "To build a user-friendly platform that seamlessly connects pet owners with trustworthy pet sitters, enhancing the overall pet care experience."
	visionStatement := flag.String("vision", vis, "The vision statement to fulfill")
    flag.Parse()

    // Call the function to validate the vision statement
    if err := validateVisionStatement(visionStatement); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

	// // Initialize the DockerClient system
	// dc, err := env.NewDockerClient(*visionStatement)
	// if err != nil{
	// 	log.Fatalln(err)
	// }
	// defer dc.Close()
	// Parse the vision statement and extract tasks/sub-goals

	vision := model.NewVision()

	// vision.Description = "To revolutionize the fitness industry by empowering users to achieve their fitness goals through personalized workout plans on a mobile application."
	vision.Description = *visionStatement


	pj := builder.NewVisionBuilder(vision.Description)

	// Define the initial vision
	fmt.Print("\nEnter your vision statement and press ENTER: \n")
	
	reader := bufio.NewReader(os.Stdin)

	inputCode, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	inputCode = strings.TrimSpace(inputCode)
	if inputCode != "" {
		vision.Description = inputCode
	}

	// Iteratively refine goals and tasks based on new insights or changes
	var (
		goals []model.Goal
		// err error
	)
	//uncheck
	_ = goals
	fmt.Printf("\nBase Vision: %s\n", vision.Description)
	// fmt.Printf("\nActiculated Vision: %s\n\n", vacticulated)
	for {
		pj.VisionEnhancement(vision)
		
	}

	// // Execute tasks to fulfill the vision statement
	// for _, task := range tasks {
	// 	// Execute the task
	// 	messages, _, done, info := dc.Rollout(task /* Pass context if needed */, true)

	// 	// Print messages and track progress
	// 	fmt.Println(messages)
	// 	if done {
	// 		// Handle completion
	// 		fmt.Println("Task completed:", task)
	// 	} else {
	// 		// Handle ongoing task
	// 		fmt.Println("Task in progress:", task)
	// 	}
	// }

	// // Generate output
	// // Include any relevant information about task completion, errors, etc.
}


