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
		filename += altFlag + ".json"
	} else {
		filename += "data.json"
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
