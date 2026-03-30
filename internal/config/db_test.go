package config

import (
	"strings"
	"testing"
)

func TestMySQLDSN_MissingEnv(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_NAME", "")
	_, err := MySQLDSN()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "faltan variables") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestMySQLDSN_OK(t *testing.T) {
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_PASSWORD", "p")
	t.Setenv("DB_NAME", "db")
	dsn, err := MySQLDSN()
	if err != nil {
		t.Fatal(err)
	}
	if dsn == "" || !strings.Contains(dsn, "localhost:3306") {
		t.Fatalf("dsn: %s", dsn)
	}
}
