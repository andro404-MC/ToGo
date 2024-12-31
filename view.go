package main

import (
	"fmt"
)

func (m model) View() string {
	s := "\n"

	switch m.state {
	case APP_VIEW:
		if len(m.TaskList) == 0 {
			s += " • empty\n"
		} else {
			for i, choice := range m.TaskList {
				cursor := "▍"
				if m.cursor == i {
					cursor = "▉"
				}

				taskLine := choice.TaskText
				if m.TaskList[i].IsSelected {
					taskLine = output.String(choice.TaskText).
						CrossOut().
						Faint().
						String()
				}

				switch m.TaskList[i].Rating {
				case 1:
					taskLine = output.String(taskLine).Foreground(output.Color("#f3bb1b")).String()
				case 2:
					taskLine = output.String(taskLine).Foreground(output.Color("#f13637")).String()
				}

				s += fmt.Sprintf("%s %s\n", cursor, taskLine)
			}
		}

	case APP_ADD:
		s += fmt.Sprintf("\n%s\n\n",
			m.textInput.View(),
		)
	}

	return s
}
