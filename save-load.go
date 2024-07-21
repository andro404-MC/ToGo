package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"gopkg.in/yaml.v3"
)

func Save(m *model) {
	yamlData, err := yaml.Marshal(&m)
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

	writenData := strings.ReplaceAll(string(yamlData), "\n    -", "\n\n    -")
	writenData = strings.ReplaceAll(writenData, ":\n", ":")
	err = os.WriteFile(filename, []byte(writenData), 0o644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func Load() model {
	ti := textinput.New()
	ti.Placeholder = "To Do"
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 35

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return model{
			TaskList: []task{
				{"Do this", true, 0},
				{"Do that", false, 1},
				{"this is a never ending cycle", false, 2},
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

	err = yaml.Unmarshal(loadedData, &loadedModel)
	if err != nil {
		fmt.Println("Error:", err)
	}

	loadedModel.textInput = ti
	loadedModel.state = 1
	return loadedModel
}
