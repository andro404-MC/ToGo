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
		fmt.Printf(
			"Buddy what have you done ?? i mean this is literaly just a task app how the fuck did you brock items Error.\nError: %v",
			err,
		)
		os.Exit(1)
	}
}

type model struct {
	choices   []string
	cursor    int
	selected  map[int]struct{}
	isTyping  bool
	textInput textinput.Model
	err       error
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
	isEmpty := false
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			Save(m)
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if m.isTyping {
				if m.textInput.Value() > "" {
					m.choices = append(m.choices, m.textInput.Value())
					m.isTyping = false
					isEmpty = true
				}
			} else {
				m.isTyping = true
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	case errMsg:
		m.err = msg
		return m, nil
	}
	if m.isTyping {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	if isEmpty {
		m.textInput.Reset()
	}
	return m, cmd
}

func (m model) View() string {
	if m.isTyping {
		return fmt.Sprintf("Add a new task\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc to quit)",
		) + "\n"
	}

	s := "Task list:\n\n"

	for i, choice := range m.choices {

		cursor := "▍"
		if m.cursor == i {
			cursor = "▉"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✓"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress KeyEsc to quit and d to delete.\n"
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
			choices:   []string{"Do this", "Do that", "this is a never ending cycle"},
			selected:  make(map[int]struct{}),
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

	// Decode JSON data into the model
	err = json.Unmarshal(loadedData, &loadedModel)
	if err != nil {
		fmt.Println("Error:", err)
		return model{}
	}
	loadedModel.textInput = ti
	return loadedModel
}

//       ██╗ ███████╗  ██████╗  ███╗   ██╗        ██████╗  ███████╗
//       ██║ ██╔════╝ ██╔═══██╗ ████╗  ██║        ██╔══██╗ ██╔════╝
//       ██║ ███████╗ ██║   ██║ ██╔██╗ ██║ █████╗ ██████╔╝ ███████╗
//  ██   ██║ ╚════██║ ██║   ██║ ██║╚██╗██║ ╚════╝ ██╔══██╗ ╚════██║
//  ╚█████╔╝ ███████║ ╚██████╔╝ ██║ ╚████║        ██████╔╝ ███████║
//   ╚════╝  ╚══════╝  ╚═════╝  ╚═╝  ╚═══╝        ╚═════╝  ╚══════╝

func (m model) MarshalJSON() ([]byte, error) {
	type ExportedModel struct {
		Choices  []string
		Cursor   int
		Selected map[int]struct{}
	}

	exportedModel := ExportedModel{
		Choices:  m.choices,
		Cursor:   m.cursor,
		Selected: m.selected,
	}

	return json.Marshal(exportedModel)
}

func (m *model) UnmarshalJSON(data []byte) error {
	var exportedModel struct {
		Choices  []string
		Cursor   int
		Selected map[int]struct{}
	}

	if err := json.Unmarshal(data, &exportedModel); err != nil {
		return err
	}

	m.choices = exportedModel.Choices
	m.cursor = exportedModel.Cursor
	m.selected = exportedModel.Selected

	return nil
}
