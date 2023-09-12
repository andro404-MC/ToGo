package main

import "fmt"

func (m model) View() string {
	s := "\n"

	switch m.state {
	case 1:
		s += "Task list:\n\n"
		if len(m.taskList) == 0 {
			s += " • empty\n"
		} else {
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
