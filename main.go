package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	termenv "github.com/muesli/termenv"
)

var (
	filename string = "/.local/share/togo/data.json"
	output   *termenv.Output
)

func initialModel() model {
	return Load()
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func main() {
	output = termenv.NewOutput(os.Stdout)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(3)
	}
	filename = homeDir + filename
	flagStuff()

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		msg := "Buddy what have you done ?? i mean this is literaly just a task app how "
		msg += "the fuck did you broke it.\nWell you dont have to gess here is the error :\n%v"
		fmt.Printf(msg, err)
		os.Exit(1)
	}
}
