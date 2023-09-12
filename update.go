package main

import tea "github.com/charmbracelet/bubbletea"

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

func addNew(m *model, newTask task) {
	m.taskList = append(m.taskList, newTask)
	m.cursor = 0
	m.state = 1
	m.textInput.Reset()
}
