package model



// Function to add a new task to the user story
func (us *UserStory) AddTask(task *Task) {
    us.Tasks = append(us.Tasks, task)
}
// Function to mark a task as completed within the user story
func (us *UserStory) MarkTaskCompleted(taskID int) Error{
    for i, task := range us.Tasks {
        if task.ID == taskID {
            us.Tasks[i].Completed = true
            return ""
        }
    }
    return TaskNotFound
}
// Function to calculate the overall progress of the user story (percentage of completed tasks)
func (us *UserStory) Progress() float64 {
    if len(us.Tasks) == 0 {
        return 0 // Avoid division by zero
    }
    completedTasks := 0
    for _, task := range us.Tasks {
        if task.Completed {
            completedTasks++
        }
    }
    return float64(completedTasks) / float64(len(us.Tasks)) * 100
}
// Function to calculate a team member's workload (number of assigned tasks)
func (us *UserStory) Workload(memberName string) int {
    workload := 0
    for _, task := range us.Tasks {
        if task.AssignedTo == memberName {
            workload++
        }
    }
    return workload
}
// Function to track a team member's progress on assigned tasks
func (us *UserStory) ProgressForMember(memberName string) float64 {
    assignedTasks := 0
    completedTasks := 0
    for _, task := range us.Tasks {
        if task.AssignedTo == memberName {
            assignedTasks++
            if task.Completed {
                completedTasks++
            }
        }
    }
    if assignedTasks == 0 {
        return 0 // Avoid division by zero
    }
    return float64(completedTasks) / float64(assignedTasks) * 100
}
// Function to notify team members when their assigned tasks are completed
func (us *UserStory) NotifyTaskCompletion(memberName string) string{
    for _, task := range us.Tasks {
        if task.AssignedTo == memberName && task.Completed {
            return "Task "+task.Description+" assigned to "+memberName+" has been completed."
        }
    }
    return ""
}
// Function to automatically balance workload by reassigning tasks among team members
func (us *UserStory) BalanceWorkload() {
    // Calculate average workload
    totalTasks := len(us.Tasks)
    averageWorkload := totalTasks / len(us.Tasks)
    
    // Track workload for each team member
    workloadMap := make(map[string]int)
    for _, task := range us.Tasks {
        workloadMap[task.AssignedTo]++
    }
    
    // Reassign tasks to balance workload
    for member, workload := range workloadMap {
        if workload > averageWorkload {
            for i := 0; i < workload-averageWorkload; i++ {
                for j, task := range us.Tasks {
                    if task.AssignedTo == member {
                        us.Tasks[j].AssignedTo = ""
                        break
                    }
                }
            }
        }
    }
}

// Define interface for team members
type TeamMember interface {
	WorkOn(userStory string)
}

