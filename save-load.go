package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
)

func Save(m *model) {
	jsonData, err := json.Marshal(m)
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

	// Write JSON data to a file
	err = os.WriteFile(filename, jsonData, 0o644)
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
			taskList: []task{
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

	err = json.Unmarshal(loadedData, &loadedModel)
	if err != nil {
		fmt.Println("Error:", err)
	}

	loadedModel.textInput = ti
	loadedModel.state = 1
	return loadedModel
}

type ExportedModel struct {
	TaskList []task
	Cursor   int
}

func (m model) MarshalJSON() ([]byte, error) {
	exportedModel := ExportedModel{
		TaskList: m.taskList,
		Cursor:   m.cursor,
	}

	return json.Marshal(exportedModel)
}

func (m *model) UnmarshalJSON(data []byte) error {
	var exportedModel ExportedModel

	if err := json.Unmarshal(data, &exportedModel); err != nil {
		return err
	}

	m.taskList = exportedModel.TaskList
	m.cursor = exportedModel.Cursor

	return nil
}

type ExportedTask struct {
	TaskText   string
	IsSelected bool
	Rating     int8
}

func (t task) MarshalJSON() ([]byte, error) {
	exportedTask := ExportedTask{
		TaskText:   t.taskText,
		IsSelected: t.isSelected,
		Rating:     t.rating,
	}

	return json.Marshal(exportedTask)
}

func (t *task) UnmarshalJSON(data []byte) error {
	var exportedTask ExportedTask

	if err := json.Unmarshal(data, &exportedTask); err != nil {
		return err
	}

	t.taskText = exportedTask.TaskText
	t.isSelected = exportedTask.IsSelected
	t.rating = exportedTask.Rating

	return nil
}
