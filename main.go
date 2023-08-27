package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	termenv "github.com/muesli/termenv"
)

var (
	filename string = "/.local/share/togo/data.json"
	output   *termenv.Output
)

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

func flagStuff() {
	m := Load()
	taskFlag := ""
	var didLs, willQuit, printHelp bool

	flag.StringVar(&taskFlag, "add", "", "Task")
	flag.BoolVar(&didLs, "ls", false, "List")
	flag.BoolVar(&printHelp, "h", false, "Help")
	flag.Parse()

	if taskFlag != "" {
		addNew(&m, task{taskFlag, false})
		fmt.Println("New task has been added : \n", taskFlag)
		willQuit = true
		Save(&m)
	}

	if didLs {
		s := "\n"
		for i, choice := range m.taskList {
			taskLine := fmt.Sprintf("%s", choice.taskText)
			if m.taskList[i].isSelected {
				taskLine = output.String(fmt.Sprintf("%s", choice.taskText)).
					CrossOut().
					Faint().
					String()
			}

			s += fmt.Sprintf(" • %s\n", taskLine)
		}
		fmt.Println(s)
		willQuit = true
	}

	if printHelp {
		flag.PrintDefaults()
		willQuit = true
	}

	if willQuit {
		os.Exit(0)
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

func addNew(m *model, newTask task) {
	m.taskList = append(m.taskList, newTask)
	m.cursor = 0
	m.state = 1
	m.textInput.Reset()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.state == 2 {
				if m.textInput.Value() > "" {
					addNew(&m, task{m.textInput.Value(), false})
				}
			} else {
				m.state = 2
			}
		case "ctrl+c":
			Save(&m)
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.state == 2 {
		m.textInput, cmd = m.textInput.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.state = 1
				m.textInput.Reset()
			}
		case errMsg:
			m.err = msg
			return m, nil
		}
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else if m.cursor == 0 {
				m.cursor = len(m.taskList) - 1
			}
		case "down":
			if m.cursor < len(m.taskList)-1 {
				m.cursor++
			} else if m.cursor == len(m.taskList)-1 {
				m.cursor = 0
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
		case "esc":
			Save(&m)
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	s := "\n"

	switch m.state {
	case 1:
		s += "Task list:\n\n"

		for i, choice := range m.taskList {

			cursor := "▍"
			if m.cursor == i {
				cursor = "▉"
			}

			taskLine := fmt.Sprintf("%s", choice.taskText)
			if m.taskList[i].isSelected {
				taskLine = output.String(fmt.Sprintf("%s", choice.taskText)).
					CrossOut().
					Faint().
					String()
			}

			s += fmt.Sprintf("%s %s\n", cursor, taskLine)
		}

		s += "\n(Esc) quit - (d) delete - (enter) add\n"

	case 2:
		s += fmt.Sprintf("Add a new task\n\n%s\n\n%s",
			m.textInput.View(),
			"(esc) return - (enter) confirm\n",
		)
	}

	return s
}

// ███████╗  █████╗  ██╗   ██╗ ███████╗     ██╗ ██╗       ██████╗   █████╗  ██████╗
// ██╔════╝ ██╔══██╗ ██║   ██║ ██╔════╝    ██╔╝ ██║      ██╔═══██╗ ██╔══██╗ ██╔══██╗
// ███████╗ ███████║ ██║   ██║ █████╗     ██╔╝  ██║      ██║   ██║ ███████║ ██║  ██║
// ╚════██║ ██╔══██║ ╚██╗ ██╔╝ ██╔══╝    ██╔╝   ██║      ██║   ██║ ██╔══██║ ██║  ██║
// ███████║ ██║  ██║  ╚████╔╝  ███████╗ ██╔╝    ███████╗ ╚██████╔╝ ██║  ██║ ██████╔╝
// ╚══════╝ ╚═╝  ╚═╝   ╚═══╝   ╚══════╝ ╚═╝     ╚══════╝  ╚═════╝  ╚═╝  ╚═╝ ╚═════╝

func Save(m *model) {
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
			state:     1,
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
