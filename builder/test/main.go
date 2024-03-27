package main

import (
	"ai_agents/vision_builder/builder"
	"ai_agents/vision_builder/model"
	"fmt"
)

// func Test_MineWorkerStructs(t *testing.T){
func main(){
output := `// PetSitter representing individuals who offer pet sitting services
type PetSitter struct {
		ID          int
		Username    string
		Email       string
		Phone       string
		Experience  int
		Rating      float64
		BackgroundCheck BackgroundCheck
}

// BackgroundCheck to manage the background check process for pet sitters    
type BackgroundCheck struct {
		ID         int
		SitterID   int
		Status     string
		Report     string
}

// Booking to handle the booking of pet sitting services
type Booking struct {
		ID           int
		PetOwner     PetOwner
		PetSitter    PetSitter
		DateTime     time.Time
		Duration     int
		Confirmed    bool
}

// Review to allow pet owners to leave feedback for pet sitters
type Review struct {
		ID          int
		PetSitter   PetSitter
		PetOwner    PetOwner
		Rating      int
		Comment     string
		CreatedAt   time.Time
}`

aa := builder.MineWorkersAndModels(output)
count := 1
for k,v := range aa{
	fmt.Println("k=",k)
	for _, parentApp :=  range v.(map[string]model.App){
		fmt.Printf("\nParent %s", parentApp.Code)
		fmt.Println(count)
		count++
		for _, child := range parentApp.Children{
			fmt.Printf("\nChild %s", child.Code)
			fmt.Println(count)
			count++
		}
	}
}

}