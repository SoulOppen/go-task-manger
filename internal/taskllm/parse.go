package taskllm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/SoulOppen/task-manager-go/internal/task"
)

var (
	// ErrEmptyTasksList indica que el modelo devolvio tasks vacio.
	ErrEmptyTasksList = errors.New("el modelo no devolvio ninguna tarea (tasks vacio)")
)

type tasksEnvelope struct {
	Tasks []taskItem `json:"tasks"`
}

type taskItem struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Relevance     int    `json:"relevance"`
	Due           string `json:"due"`
	DependsOnID   string `json:"depends_on_id"`
}

// ParseTasksJSON decodifica la respuesta del modelo con campos desconocidos rechazados.
func ParseTasksJSON(raw string) ([]taskItem, error) {
	raw = strings.TrimSpace(raw)
	raw = stripMarkdownFences(raw)
	dec := json.NewDecoder(bytes.NewReader([]byte(raw)))
	dec.DisallowUnknownFields()
	var env tasksEnvelope
	if err := dec.Decode(&env); err != nil {
		return nil, fmt.Errorf("JSON invalido o con campos no permitidos: %w", err)
	}
	if len(env.Tasks) == 0 {
		return nil, ErrEmptyTasksList
	}
	return env.Tasks, nil
}

func stripMarkdownFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		// quita primera linea ``` o ```json
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		if end := strings.LastIndex(s, "```"); end > 0 {
			s = strings.TrimSpace(s[:end])
		}
	}
	return strings.TrimSpace(s)
}

// BuildTasks convierte items en *task.Task listos para Validate y Create.
func BuildTasks(items []taskItem) ([]*task.Task, error) {
	out := make([]*task.Task, 0, len(items))
	for i, it := range items {
		rel := it.Relevance
		if rel < 1 || rel > 10 {
			rel = 5
		}
		due, err := task.ParseDueDate(it.Due)
		if err != nil {
			return nil, fmt.Errorf("tarea %d: %w", i+1, err)
		}
		t := task.NewTask(strings.TrimSpace(it.Name), strings.TrimSpace(it.Description), rel, due)
		if dep := strings.TrimSpace(it.DependsOnID); dep != "" {
			t.DependsOnID = &dep
		}
		if err := t.Validate(); err != nil {
			return nil, fmt.Errorf("tarea %d: %w", i+1, err)
		}
		out = append(out, t)
	}
	return out, nil
}
