package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case 1:
		// List
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j", "down":
				if m.cursor < len(m.taskList)-1 {
					m.cursor++
				}
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "G":
				m.cursor = len(m.taskList) - 1
			case "g":
				m.cursor = 0
			case " ":
				m.taskList[m.cursor].isSelected = !m.taskList[m.cursor].isSelected
			case "q":
				Save(&m)
				return m, tea.Quit
			case "a":
				m.state = 2
			case "d":
				if len(m.taskList) > 0 {
					m.taskList = append(m.taskList[:m.cursor], m.taskList[m.cursor+1:]...)
					if m.cursor > len(m.taskList)-1 {
						m.cursor = len(m.taskList) - 1
					}
				}
			}
		case errMsg:
			m.err = msg
			return m, nil
		}
	case 2:
		// add new element
		m.textInput, cmd = m.textInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.textInput.Value() > "" {
					addNew(&m, task{m.textInput.Value(), false, 0})
					m.textInput.Reset()
					Save(&m)
					m.state = 1
				}
			case "esc":
				m.state = 1
				m.textInput.Reset()
			}
		case errMsg:
			m.err = msg
			return m, nil
		}
	}

	// HACK
	// ALL states
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctr+c":
			Save(&m)
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, cmd
}

func addNew(m *model, newTask task) {
	m.taskList = append(m.taskList, newTask)
	m.cursor = 0
	m.state = 1
}
