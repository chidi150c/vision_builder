package main

import (
	"ai_agents/vision_builder/builder"
	"ai_agents/vision_builder/model"
	"bufio"
	// "context"
	"fmt"
	"log"
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
		_ = pj.Learn(true)
		pj.VisionStatement = vision.Description
		old := ""
		for k,v := range vision.Goals{
			if old == v.Description{
				delete(vision.Goals, k)
				continue
			}
			old = v.Description
			fmt.Println(k,".", v.Description)
			Tasks := pj.DeriveTasksFromGoal(v, pj.VisionStatement)
			
			// Execute tasks to fulfill the vision statement
			for _, task := range Tasks {
				
				// Execute the task
				messages, _, done, _ := pj.Rollout(task.Description, pj.Context, true)

				// Print messages and track progress
				fmt.Println(messages)
				if done {
					// Handle completion
					fmt.Println("Task completed:", task)
				} else {
					// Handle ongoing task
					fmt.Println("Task in progress:", task)
				}
			}
		}
	}


	// // Generate output
	// // Include any relevant information about task completion, errors, etc.
}


