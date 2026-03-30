package config

import (
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultName    = "task-manager-go"
	DefaultVersion = "0.0.1"

	DefaultShortDescription = "CLI para gestion de tareas"
	DefaultLongDescription  = "Herramienta CLI para administrar tareas, autenticacion y flujo diario desde terminal."
)

func init() {
	// Carga .env si existe (p. ej. credenciales DB); nombre y version no vienen de env.
	_ = godotenv.Load()
}

func AppName() string {
	return strings.ReplaceAll(strings.TrimSpace(DefaultName), " ", "-")
}

func AppVersion() string {
	return strings.TrimSpace(DefaultVersion)
}

func AppShortDescription() string {
	return DefaultShortDescription
}

func AppLongDescription() string {
	return DefaultLongDescription
}
