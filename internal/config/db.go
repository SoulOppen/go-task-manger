package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

// MySQLDSN construye el DSN para database/sql + mysql driver.
// Requiere DB_HOST, DB_PORT, DB_USER, DB_NAME. DB_PASSWORD puede estar vacío en desarrollo.
func MySQLDSN() (string, error) {
	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	pass := os.Getenv("DB_PASSWORD")
	name := strings.TrimSpace(os.Getenv("DB_NAME"))

	if host == "" || port == "" || user == "" || name == "" {
		return "", fmt.Errorf("faltan variables DB_HOST, DB_PORT, DB_USER o DB_NAME en el entorno (.env)")
	}

	cfg := mysqldriver.NewConfig()
	cfg.User = user
	cfg.Passwd = pass
	cfg.Net = "tcp"
	cfg.Addr = host + ":" + port
	cfg.DBName = name
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	return cfg.FormatDSN(), nil
}
