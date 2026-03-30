package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMigrate_RunsDDL(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS tasks").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS users").WillReturnResult(sqlmock.NewResult(0, 0))

	if err := Migrate(context.Background(), database); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMigrate_tasksExecError(t *testing.T) {
	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS tasks").WillReturnError(errors.New("ddl fallo"))

	if err := Migrate(context.Background(), database); err == nil {
		t.Fatal("expected error")
	}
}

func TestWithDB_SinVariablesEntorno(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")
	err := WithDB(context.Background(), func(db *sql.DB) error {
		t.Fatal("fn no debe ejecutarse sin DSN")
		return nil
	})
	if err == nil {
		t.Fatal("expected error from MySQLDSN / Open")
	}
}
