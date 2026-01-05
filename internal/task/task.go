package task

import (
	"fmt"
)

type Task struct {
	ID      string  `json:"id"`
	Proyect *string `json:"project,omitempty"`
	Name    string  `json:"name"`
}

var agenda []Task

func addTask(task Task) {
	agenda = append(agenda, task)
	for _, t := range agenda {
		t.printTask()
	}
}
func (t Task) printTask() {
	if t.Project != nil {
		fmt.Printf("%s - %s - %s\n", t.ID, *t.Project, t.Name)
	} else {
		fmt.Printf("%s - (sin proyecto) - %s\n", t.ID, t.Name)
	}
}
