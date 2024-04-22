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
	rating     int8
}

// rating
// 0 normal
// 1 important
// 2 critical

type (
	errMsg error
)
