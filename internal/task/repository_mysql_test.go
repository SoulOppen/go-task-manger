package task

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)
	created := time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC)
	tt := &Task{
		ID:          "550e8400-e29b-41d4-a716-446655440000",
		Name:        "uno",
		Description: "desc",
		Relevance:   5,
		CreatedAt:   created,
		DueDate:     nil,
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(tt.ID, tt.Name, tt.Description, tt.Relevance, tt.CreatedAt, nil).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.Create(context.Background(), tt); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_ListOrdered(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ct := time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "relevance", "created_at", "due_date"}).
		AddRow("550e8400-e29b-41d4-a716-446655440000", "n", "d", 8, ct, nil)

	mock.ExpectQuery("SELECT (.+) FROM tasks ORDER BY relevance DESC").
		WillReturnRows(rows)

	list, err := repo.ListOrdered(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Fatalf("unexpected list: %#v", list)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
