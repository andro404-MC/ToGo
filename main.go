package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var filename string = "/.local/share/tuiapptest/data.json"

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(3)
	}
	filename = homeDir + filename

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		msg := "Buddy what have you done ?? i mean this is literaly just a task app how "
		msg += "the fuck did you brock items Error.\nError: %v"
		fmt.Printf(msg, err)
		os.Exit(1)
	}
}

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

func initialModel() model {
	return Load()
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			Save(m)
			return m, tea.Quit
		case "enter":
			if m.state == 2 {
				if m.textInput.Value() > "" {
					m.taskList = append(m.taskList, task{m.textInput.Value(), false})
					m.state = 1
					m.textInput.Reset()
				}
			} else {
				m.state = 2
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.state == 2 {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.taskList)-1 {
				m.cursor++
			}
		case "d":
			if len(m.taskList) > 0 {
				m.taskList = append(m.taskList[:m.cursor], m.taskList[m.cursor+1:]...)
				if m.cursor > len(m.taskList)-1 {
					m.cursor = len(m.taskList) - 1
				}
			}
		case " ":
			m.taskList[m.cursor].isSelected = !m.taskList[m.cursor].isSelected
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	switch m.state {
	case 1:
		s = "Task list:\n\n"

		for i, choice := range m.taskList {

			cursor := "▍"
			if m.cursor == i {
				cursor = "▉"
			}

			checked := " "
			if m.taskList[i].isSelected {
				checked = "✓"
			}

			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.taskText)
		}

		s += "\nPress KeyEsc to quit and d to delete.\n"

	case 2:
		s = fmt.Sprintf("Add a new task\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}

	return s
}

func Save(m model) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dirPath := filepath.Dir(filename)
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Write JSON data to a file
	err = os.WriteFile(filename, jsonData, 0o644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func Load() model {
	ti := textinput.New()
	ti.Placeholder = "To Do"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return model{
			taskList: []task{
				{"Do this", true},
				{"Do that", false},
				{"this is a never ending cycle", false},
			},
			textInput: ti,
			err:       nil,
		}
	}

	loadedData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return model{}
	}

	var loadedModel model

	err = json.Unmarshal(loadedData, &loadedModel)
	if err != nil {
		fmt.Println("Error:", err)
		return model{}
	}

	loadedModel.textInput = ti
	loadedModel.state = 1
	return loadedModel
}

//       ██╗ ███████╗  ██████╗  ███╗   ██╗        ██████╗  ███████╗
//       ██║ ██╔════╝ ██╔═══██╗ ████╗  ██║        ██╔══██╗ ██╔════╝
//       ██║ ███████╗ ██║   ██║ ██╔██╗ ██║ █████╗ ██████╔╝ ███████╗
//  ██   ██║ ╚════██║ ██║   ██║ ██║╚██╗██║ ╚════╝ ██╔══██╗ ╚════██║
//  ╚█████╔╝ ███████║ ╚██████╔╝ ██║ ╚████║        ██████╔╝ ███████║
//   ╚════╝  ╚══════╝  ╚═════╝  ╚═╝  ╚═══╝        ╚═════╝  ╚══════╝

type ExportedModel struct {
	TaskList []task
	Cursor   int
}

func (m model) MarshalJSON() ([]byte, error) {
	exportedModel := ExportedModel{
		TaskList: m.taskList,
		Cursor:   m.cursor,
	}

	return json.Marshal(exportedModel)
}

func (m *model) UnmarshalJSON(data []byte) error {
	var exportedModel ExportedModel

	if err := json.Unmarshal(data, &exportedModel); err != nil {
		return err
	}

	m.taskList = exportedModel.TaskList
	m.cursor = exportedModel.Cursor

	return nil
}

type ExportedTask struct {
	TaskText   string
	IsSelected bool
}

func (t task) MarshalJSON() ([]byte, error) {
	exportedTask := ExportedTask{
		TaskText:   t.taskText,
		IsSelected: t.isSelected,
	}

	return json.Marshal(exportedTask)
}

func (t *task) UnmarshalJSON(data []byte) error {
	var exportedTask ExportedTask

	if err := json.Unmarshal(data, &exportedTask); err != nil {
		return err
	}

	t.taskText = exportedTask.TaskText
	t.isSelected = exportedTask.IsSelected

	return nil
}
