package model

// General errors.
const (
	TaskNotFound = Error("Task Not Found")
)


// Error represents a User error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }