package main

import "github.com/charmbracelet/bubbles/textinput"

type model struct {
	taskList  []task
	cursor    int
	state     int
	textInput textinput.Model
	err       error
}

type task struct {
	taskText   string
	isSelected bool
}

type (
	errMsg error
)
