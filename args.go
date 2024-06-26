package main

import (
	"flag"
	"fmt"
	"os"
)

func flagStuff() {
	var taskFlag, altFlag string
	var didLs, willQuit, printHelp bool

	flag.StringVar(&altFlag, "t", "", "Custom list")
	flag.StringVar(&taskFlag, "a", "", "Add yask")
	flag.BoolVar(&didLs, "l", false, "List tasks")
	flag.BoolVar(&printHelp, "h", false, "Help")
	flag.Parse()

	if altFlag != "" {
		filename += altFlag + ".yaml"
	} else {
		filename += "data.yaml"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(3)
	}
	filename = homeDir + filename

	m := Load()

	if taskFlag != "" {
		addNew(&m, task{taskFlag, false, 0})
		fmt.Println("New task has been added : \n", taskFlag)
		willQuit = true
		Save(&m)
	}

	if didLs {
		s := "\n"
		for i, choice := range m.TaskList {
			taskLine := fmt.Sprintf("%s", choice.TaskText)
			if m.TaskList[i].IsSelected {
				taskLine = output.String(fmt.Sprintf("%s", choice.TaskText)).
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
