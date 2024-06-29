package main

import "github.com/charmbracelet/bubbles/textinput"

type model struct {
	TaskList []task `yaml:"tasks"`

	termHeight int
	cursor     int
	state      int

	textInput textinput.Model
	err       error
}

type task struct {
	TaskText   string `yaml:"name"`
	IsSelected bool   `yaml:"done"`
	Rating     int8   `yaml:"rating"`
}

// rating
// 0 normal
// 1 important
// 2 critical

type (
	errMsg error
)
