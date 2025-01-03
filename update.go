package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case APP_VIEW:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "n":
				if m.cursor < len(m.TaskList)-1 {
					currentTask := m.TaskList[m.cursor]
					replacedTask := m.TaskList[m.cursor+1]

					m.TaskList[m.cursor] = replacedTask
					m.TaskList[m.cursor+1] = currentTask

					m.cursor++
				}
			case "m":
				if m.cursor > 0 {
					currentTask := m.TaskList[m.cursor]
					replacedTask := m.TaskList[m.cursor-1]

					m.TaskList[m.cursor] = replacedTask
					m.TaskList[m.cursor-1] = currentTask

					m.cursor--
				}
			case "j", "down":
				if m.cursor < len(m.TaskList)-1 {
					m.cursor++
				}
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "G", "pgdown":
				m.cursor = len(m.TaskList) - 1
			case "g", "pgup":
				m.cursor = 0
			case " ":
				m.TaskList[m.cursor].IsSelected = !m.TaskList[m.cursor].IsSelected
			case "q":
				Save(&m)
				return m, tea.Quit
			case "a":
				m.state = APP_ADD
			case "d":
				if len(m.TaskList) > 0 {
					m.TaskList = append(m.TaskList[:m.cursor], m.TaskList[m.cursor+1:]...)
					if m.cursor > len(m.TaskList)-1 {
						m.cursor = len(m.TaskList) - 1
					}
				}
			case "0":
				m.TaskList[m.cursor].Rating = 0
			case "1":
				if m.TaskList[m.cursor].Rating == 1 {
					m.TaskList[m.cursor].Rating = 0
				} else {
					m.TaskList[m.cursor].Rating = 1
				}
			case "2":
				if m.TaskList[m.cursor].Rating == 2 {
					m.TaskList[m.cursor].Rating = 0
				} else {
					m.TaskList[m.cursor].Rating = 2
				}
			}
		case errMsg:
			m.err = msg
			return m, nil
		}
	case APP_ADD:
		m.textInput, cmd = m.textInput.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.textInput.Value() > "" {
					addNew(&m, task{m.textInput.Value(), false, 0})
					m.textInput.Reset()
					Save(&m)
					m.state = APP_VIEW
				}
			case "esc":
				m.state = APP_VIEW
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

	m.termHeight = 88
	return m, cmd
}

func addNew(m *model, newTask task) {
	m.TaskList = append(m.TaskList, newTask)
	m.cursor = 0
	m.state = APP_VIEW
}
