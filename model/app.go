package model

type App struct{
	Name string
	Description string
	Code string
	Children map[string]App
}

