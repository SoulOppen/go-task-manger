package task

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidate_NewTask_OK(t *testing.T) {
	d := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	tt := NewTask("Comprar", "Leche", 7, &d)
	if err := tt.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestParseDueDate(t *testing.T) {
	t.Run("vacio", func(t *testing.T) {
		d, err := ParseDueDate("  ")
		if err != nil {
			t.Fatal(err)
		}
		if d != nil {
			t.Fatal("expected nil")
		}
	})
	t.Run("ok", func(t *testing.T) {
		d, err := ParseDueDate("2026-06-15")
		if err != nil {
			t.Fatal(err)
		}
		if d == nil || d.Year() != 2026 || d.Month() != 6 || d.Day() != 15 {
			t.Fatalf("got %v", d)
		}
	})
	t.Run("invalido", func(t *testing.T) {
		_, err := ParseDueDate("15-06-2026")
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestValidate_nil(t *testing.T) {
	var tt *Task
	if err := tt.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_uuidInvalido(t *testing.T) {
	tt := NewTask("a", "b", 5, nil)
	tt.ID = "no-es-uuid"
	if err := tt.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_dependsOnSelf(t *testing.T) {
	tt := NewTask("a", "b", 5, nil)
	id := tt.ID
	tt.DependsOnID = &id
	if err := tt.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_dependsOnInvalidUUID(t *testing.T) {
	tt := NewTask("a", "b", 5, nil)
	bad := "no-uuid"
	tt.DependsOnID = &bad
	if err := tt.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestValidate_errors(t *testing.T) {
	cases := []struct {
		name string
		task func() *Task
	}{
		{"nombre vacio", func() *Task {
			t := NewTask("", "desc", 5, nil)
			t.ID = uuid.NewString()
			return t
		}},
		{"descripcion vacia", func() *Task {
			t := NewTask("n", "", 5, nil)
			t.ID = uuid.NewString()
			return t
		}},
		{"relevancia baja", func() *Task {
			t := NewTask("n", "d", 0, nil)
			t.ID = uuid.NewString()
			return t
		}},
		{"relevancia alta", func() *Task {
			t := NewTask("n", "d", 11, nil)
			t.ID = uuid.NewString()
			return t
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.task().Validate(); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestSortTasksForList_vacio(t *testing.T) {
	if got := SortTasksForList(nil); len(got) != 0 {
		t.Fatalf("len=%d", len(got))
	}
	if got := SortTasksForList([]Task{}); len(got) != 0 {
		t.Fatalf("len=%d", len(got))
	}
}

func TestSortTasksForList(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	t1 := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	dNear := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	dFar := time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC)

	tasks := []Task{
		{ID: "a", Name: "low-rel", Description: "x", Relevance: 5, CreatedAt: t0, DueDate: &dNear},
		{ID: "b", Name: "high-later", Description: "x", Relevance: 9, CreatedAt: t0, DueDate: &dFar},
		{ID: "c", Name: "high-sooner", Description: "x", Relevance: 9, CreatedAt: t1, DueDate: &dNear},
		{ID: "d", Name: "no-due", Description: "x", Relevance: 9, CreatedAt: t1, DueDate: nil},
	}

	got := SortTasksForList(tasks)
	ids := make([]string, len(got))
	for i, x := range got {
		ids[i] = x.ID
	}
	want := []string{"c", "b", "d", "a"}
	if len(ids) != len(want) {
		t.Fatalf("len=%d", len(ids))
	}
	for i := range want {
		if ids[i] != want[i] {
			t.Fatalf("position %d: got %v want %v", i, ids, want)
		}
	}
}
