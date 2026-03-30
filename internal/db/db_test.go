package db

import (
	"context"
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
