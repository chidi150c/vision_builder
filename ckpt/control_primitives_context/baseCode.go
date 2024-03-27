package main

type User struct{

}
type UserService struct{

}
type Worker interface{

}
var _ Worker = &UserService{}

func(u *UserService) GoalConcept(user User){
	//implement Tasks and Subtasks
}

func main(){

}