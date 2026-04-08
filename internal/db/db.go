package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SoulOppen/task-manager-go/internal/config"

	mysqldriver "github.com/go-sql-driver/mysql"

	_ "github.com/go-sql-driver/mysql" // registro del driver "mysql"
)

const migrateTasksSQL = `CREATE TABLE IF NOT EXISTS tasks (
  id CHAR(36) NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  relevance TINYINT NOT NULL,
  created_at DATETIME(6) NOT NULL,
  due_date DATE NULL,
  depends_on_id CHAR(36) NULL,
  CONSTRAINT fk_tasks_depends_on FOREIGN KEY (depends_on_id) REFERENCES tasks(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

const migrateUsersSQL = `CREATE TABLE IF NOT EXISTS users (
  username VARCHAR(255) NOT NULL PRIMARY KEY,
  password_hash VARCHAR(255) NOT NULL,
  quick_connect_value VARCHAR(32) NULL,
  quick_connect_created_at DATETIME(6) NULL,
  quick_connect_reset_date DATE NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

// Migrate crea tablas tasks y users si no existen y aplica migraciones incrementales idempotentes.
func Migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, migrateTasksSQL); err != nil {
		return fmt.Errorf("tasks: %w", err)
	}
	if err := migrateTasksDependsOn(ctx, db); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, migrateUsersSQL); err != nil {
		return fmt.Errorf("users: %w", err)
	}
	return nil
}

func isMySQLDuplicateDDL(err error) bool {
	var me *mysqldriver.MySQLError
	if errors.As(err, &me) {
		switch me.Number {
		case 1060, 1061, 1826: // columna/índice/FK duplicado
			return true
		}
	}
	return false
}

// migrateTasksDependsOn añade depends_on_id y FK en bases creadas antes de esa columna (CREATE IF NOT EXISTS no actualiza esquema).
func migrateTasksDependsOn(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `ALTER TABLE tasks ADD COLUMN depends_on_id CHAR(36) NULL`); err != nil && !isMySQLDuplicateDDL(err) {
		return fmt.Errorf("tasks depends_on_id: %w", err)
	}
	if _, err := db.ExecContext(ctx,
		`ALTER TABLE tasks ADD CONSTRAINT fk_tasks_depends_on FOREIGN KEY (depends_on_id) REFERENCES tasks(id) ON DELETE SET NULL`,
	); err != nil && !isMySQLDuplicateDDL(err) {
		return fmt.Errorf("tasks fk depends_on: %w", err)
	}
	return nil
}

// Open abre MySQL, hace ping y Migrate.
func Open(ctx context.Context) (*sql.DB, error) {
	dsn, err := config.MySQLDSN()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir mysql: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	if err := Migrate(pingCtx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrar: %w", err)
	}
	return db, nil
}

// WithDB abre la base, ejecuta fn y cierra la conexion.
func WithDB(ctx context.Context, fn func(*sql.DB) error) error {
	db, err := Open(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	return fn(db)
}
