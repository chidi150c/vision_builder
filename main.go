package main

import (
	"ai_agents/vision_builder/builder"
	"ai_agents/vision_builder/model"
	// "ai_agents/vision_builder/prompts"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	vis := "To build a user-friendly platform that seamlessly connects pet owners with trustworthy pet sitters, enhancing the overall pet care experience."
	reader := bufio.NewReader(os.Stdin)
	_ = vis
	pj := builder.NewVisionBuilder(reader)
	defer pj.Env.Close()
	//Get the vision
	vision := model.NewVision()
	vision.Description = vis
	// Define the initial vision
	fmt.Print("\nEnter your vision statement and press ENTER: \n")
	
	inputCode, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	inputCode = strings.TrimSpace(inputCode)
	if inputCode != "" {
		vision.Description = inputCode
	}
	pj.VisionEnhancement(vision)
	// pj.Rollout()
	pj.Learn(vision, true)
}