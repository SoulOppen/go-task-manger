package task

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const DateLayout = "2006-01-02"

// Task es el modelo persistido en MySQL.
type Task struct {
	ID          string
	Name        string
	Description string
	Relevance   int
	CreatedAt   time.Time
	DueDate     *time.Time
}

// NewTask crea una tarea nueva con id UUID y CreatedAt en UTC (no persiste).
func NewTask(name, description string, relevance int, dueDate *time.Time) *Task {
	due := normalizeDueDate(dueDate)
	return &Task{
		ID:          uuid.NewString(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Relevance:   relevance,
		CreatedAt:   time.Now().UTC(),
		DueDate:     due,
	}
}

func normalizeDueDate(d *time.Time) *time.Time {
	if d == nil {
		return nil
	}
	y, m, day := d.Date()
	t := time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
	return &t
}

// Validate valida reglas de negocio antes de crear/actualizar.
func (t *Task) Validate() error {
	if t == nil {
		return errors.New("tarea nil")
	}
	if t.Name == "" {
		return errors.New("el nombre es obligatorio")
	}
	if t.Description == "" {
		return errors.New("la descripcion es obligatoria")
	}
	if t.Relevance < 1 || t.Relevance > 10 {
		return fmt.Errorf("la relevancia debe estar entre 1 y 10; recibido %d", t.Relevance)
	}
	if _, err := uuid.Parse(t.ID); err != nil {
		return fmt.Errorf("id invalido: %w", err)
	}
	return nil
}

// ParseDueDate parsea YYYY-MM-DD o devuelve error.
func ParseDueDate(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	d, err := time.ParseInLocation(DateLayout, s, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("fecha de entrega invalida (use YYYY-MM-DD): %w", err)
	}
	return normalizeDueDate(&d), nil
}

// SortTasksForList ordena en memoria con la misma logica que ListOrdered en SQL.
func SortTasksForList(tasks []Task) []Task {
	out := make([]Task, len(tasks))
	copy(out, tasks)
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Relevance != b.Relevance {
			return a.Relevance > b.Relevance
		}
		aDue := a.DueDate != nil
		bDue := b.DueDate != nil
		if aDue != bDue {
			return aDue && !bDue
		}
		if aDue && bDue {
			ai := a.DueDate.UTC()
			bi := b.DueDate.UTC()
			if !ai.Equal(bi) {
				return ai.Before(bi)
			}
		}
		return a.CreatedAt.UTC().Before(b.CreatedAt.UTC())
	})
	return out
}
