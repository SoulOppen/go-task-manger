package auth

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mysqldriver "github.com/go-sql-driver/mysql"
)

func TestMySQLUserStore_CreateAndGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewMySQLUserStore(db)
	ctx := context.Background()

	u := User{Username: "a", PasswordHash: "hash"}
	mock.ExpectExec("INSERT INTO users").
		WithArgs(u.Username, u.PasswordHash).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := store.Create(ctx, u); err != nil {
		t.Fatal(err)
	}

	rows := sqlmock.NewRows([]string{"username", "password_hash", "quick_connect_value", "quick_connect_created_at", "quick_connect_reset_date"}).
		AddRow("a", "hash", nil, nil, nil)
	mock.ExpectQuery("SELECT .+ FROM users WHERE username").WithArgs("a").WillReturnRows(rows)

	got, err := store.GetByUsername(ctx, "a")
	if err != nil {
		t.Fatal(err)
	}
	if got.Username != "a" || got.PasswordHash != "hash" {
		t.Fatalf("%+v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestMySQLUserStore_CreateDuplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewMySQLUserStore(db)

	mock.ExpectExec("INSERT INTO users").WillReturnError(&mysqldriver.MySQLError{Number: 1062})
	if err := store.Create(context.Background(), User{Username: "x", PasswordHash: "y"}); err != ErrUserExists {
		t.Fatalf("got %v", err)
	}
}

func TestMySQLUserStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	store := NewMySQLUserStore(db)
	ctx := context.Background()

	created := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)
	u := User{
		Username:              "a",
		PasswordHash:          "hash",
		QuickConnectValue:     "abc123",
		QuickConnectCreatedAt: created.UTC().Format(time.RFC3339),
		QuickConnectResetDate: "2026-03-30",
	}

	mock.ExpectExec("UPDATE users SET").WillReturnResult(sqlmock.NewResult(0, 1))
	if err := store.Update(ctx, u); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
