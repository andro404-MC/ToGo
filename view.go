package main

import (
	"fmt"
)

func (m model) View() string {
	s := "\n"

	switch m.state {
	case 1:
		if len(m.taskList) == 0 {
			s += " • empty\n"
		} else {
			for i, choice := range m.taskList {
				cursor := "▍"
				if m.cursor == i {
					cursor = "▉"
				}

				taskLine := choice.taskText
				if m.taskList[i].isSelected {
					taskLine = output.String(choice.taskText).
						CrossOut().
						Faint().
						String()
				}

				switch m.taskList[i].rating {
				case 1:
					taskLine = output.String(taskLine).Foreground(output.Color("#f3bb1b")).String()
				case 2:
					taskLine = output.String(taskLine).Foreground(output.Color("#f13637")).String()
				}

				s += fmt.Sprintf("%s %s\n", cursor, taskLine)
			}
		}

	case 2:
		s += fmt.Sprintf("\n%s\n\n",
			m.textInput.View(),
		)
	}

	return s
}
