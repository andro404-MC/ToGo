package main

import (
	"flag"
	"fmt"
	"os"
)

func flagStuff() {
	m := Load()
	taskFlag := ""
	var didLs, willQuit, printHelp bool

	flag.StringVar(&taskFlag, "add", "", "Task")
	flag.BoolVar(&didLs, "ls", false, "List")
	flag.BoolVar(&printHelp, "h", false, "Help")
	flag.Parse()

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
