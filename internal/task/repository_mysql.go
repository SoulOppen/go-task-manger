package task

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNotFound = errors.New("tarea no encontrada")

// Repository accede a tareas en MySQL.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func dueSQLArg(d *time.Time) interface{} {
	if d == nil {
		return nil
	}
	dd := normalizeDueDate(d)
	if dd == nil {
		return nil
	}
	return dd.Format("2006-01-02")
}

func (r *Repository) Create(ctx context.Context, t *Task) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, name, description, relevance, created_at, due_date) VALUES (?, ?, ?, ?, ?, ?)`,
		t.ID, t.Name, t.Description, t.Relevance, t.CreatedAt, dueSQLArg(t.DueDate),
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Task, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, relevance, created_at, due_date FROM tasks WHERE id = ?`, id)
	t, err := scanTaskRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *Repository) Update(ctx context.Context, t *Task) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET name = ?, description = ?, relevance = ?, due_date = ? WHERE id = ?`,
		t.Name, t.Description, t.Relevance, dueSQLArg(t.DueDate), t.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListOrdered(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, description, relevance, created_at, due_date FROM tasks
		 ORDER BY relevance DESC, (due_date IS NULL), due_date ASC, created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		t, err := scanTaskRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, rows.Err()
}

func scanTaskRow(row *sql.Row) (*Task, error) {
	var t Task
	var due sql.NullTime
	if err := row.Scan(&t.ID, &t.Name, &t.Description, &t.Relevance, &t.CreatedAt, &due); err != nil {
		return nil, err
	}
	if due.Valid {
		d := normalizeDueDate(&due.Time)
		t.DueDate = d
	}
	return &t, nil
}

func scanTaskRows(rows *sql.Rows) (*Task, error) {
	var t Task
	var due sql.NullTime
	if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Relevance, &t.CreatedAt, &due); err != nil {
		return nil, err
	}
	if due.Valid {
		t.DueDate = normalizeDueDate(&due.Time)
	}
	return &t, nil
}
