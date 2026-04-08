package task

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	ErrNotFound           = errors.New("tarea no encontrada")
	ErrNoTasks            = errors.New("no hay tareas")
	ErrCircularDependency = errors.New("dependencia circular entre tareas")
	ErrDependsNotFound    = errors.New("la tarea dependiente no existe")
)

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

func dependsOnSQLArg(id *string) interface{} {
	if id == nil {
		return nil
	}
	s := *id
	if s == "" {
		return nil
	}
	return s
}

func (r *Repository) Create(ctx context.Context, t *Task) error {
	if err := r.ensureDependsOnValid(ctx, t.ID, t.DependsOnID); err != nil {
		return err
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, name, description, relevance, created_at, due_date, depends_on_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.Name, t.Description, t.Relevance, t.CreatedAt, dueSQLArg(t.DueDate), dependsOnSQLArg(t.DependsOnID),
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Task, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT t.id, t.name, t.description, t.relevance, t.created_at, t.due_date, t.depends_on_id, p.name
		 FROM tasks t
		 LEFT JOIN tasks p ON t.depends_on_id = p.id
		 WHERE t.id = ?`, id)
	t, err := scanTaskRowWithParent(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *Repository) Update(ctx context.Context, t *Task) error {
	if err := r.ensureDependsOnValid(ctx, t.ID, t.DependsOnID); err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET name = ?, description = ?, relevance = ?, due_date = ?, depends_on_id = ? WHERE id = ?`,
		t.Name, t.Description, t.Relevance, dueSQLArg(t.DueDate), dependsOnSQLArg(t.DependsOnID), t.ID,
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
		`SELECT t.id, t.name, t.description, t.relevance, t.created_at, t.due_date, t.depends_on_id, p.name
		 FROM tasks t
		 LEFT JOIN tasks p ON t.depends_on_id = p.id
		 ORDER BY t.relevance DESC, (t.due_date IS NULL), t.due_date ASC, t.created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		t, err := scanTaskRowsWithParent(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, rows.Err()
}

// PickRandom devuelve una tarea elegida al azar.
func (r *Repository) PickRandom(ctx context.Context) (*Task, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT t.id, t.name, t.description, t.relevance, t.created_at, t.due_date, t.depends_on_id, p.name
		 FROM tasks t
		 LEFT JOIN tasks p ON t.depends_on_id = p.id
		 ORDER BY RAND() LIMIT 1`)
	t, err := scanTaskRowWithParent(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoTasks
		}
		return nil, err
	}
	return t, nil
}

// ensureDependsOnValid comprueba existencia de la tarea padre y ausencia de ciclos.
func (r *Repository) ensureDependsOnValid(ctx context.Context, taskID string, dep *string) error {
	if dep == nil {
		return nil
	}
	depID := strings.TrimSpace(*dep)
	if depID == "" {
		return nil
	}
	if depID == taskID {
		return ErrCircularDependency
	}
	exists, err := r.rowExists(ctx, depID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrDependsNotFound
	}
	return r.walkDependsChain(ctx, taskID, depID)
}

func (r *Repository) rowExists(ctx context.Context, id string) (bool, error) {
	var one int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM tasks WHERE id = ? LIMIT 1`, id).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// walkDependsChain recorre la cadena de dependencias desde proposedDep; si alcanza taskID, habria ciclo.
func (r *Repository) walkDependsChain(ctx context.Context, taskID, proposedDep string) error {
	current := proposedDep
	seen := make(map[string]struct{})
	for i := 0; i < 64; i++ {
		if current == taskID {
			return ErrCircularDependency
		}
		if _, dup := seen[current]; dup {
			return ErrCircularDependency
		}
		seen[current] = struct{}{}

		var next sql.NullString
		err := r.db.QueryRowContext(ctx, `SELECT depends_on_id FROM tasks WHERE id = ?`, current).Scan(&next)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrDependsNotFound
			}
			return err
		}
		if !next.Valid || strings.TrimSpace(next.String) == "" {
			return nil
		}
		current = next.String
	}
	return ErrCircularDependency
}

func scanTaskRowWithParent(row *sql.Row) (*Task, error) {
	var t Task
	var due sql.NullTime
	var depID sql.NullString
	var depName sql.NullString
	if err := row.Scan(&t.ID, &t.Name, &t.Description, &t.Relevance, &t.CreatedAt, &due, &depID, &depName); err != nil {
		return nil, err
	}
	if due.Valid {
		d := normalizeDueDate(&due.Time)
		t.DueDate = d
	}
	if depID.Valid && depID.String != "" {
		s := depID.String
		t.DependsOnID = &s
	}
	if depName.Valid {
		t.DependsOnName = depName.String
	}
	return &t, nil
}

func scanTaskRowsWithParent(rows *sql.Rows) (*Task, error) {
	var t Task
	var due sql.NullTime
	var depID sql.NullString
	var depName sql.NullString
	if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Relevance, &t.CreatedAt, &due, &depID, &depName); err != nil {
		return nil, err
	}
	if due.Valid {
		t.DueDate = normalizeDueDate(&due.Time)
	}
	if depID.Valid && depID.String != "" {
		s := depID.String
		t.DependsOnID = &s
	}
	if depName.Valid {
		t.DependsOnName = depName.String
	}
	return &t, nil
}
