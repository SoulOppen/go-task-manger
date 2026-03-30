package config

import (
	"strings"

	"github.com/joho/godotenv"
)

const (
	DefaultName    = "gtm"
	DefaultVersion = "0.1.0"

	DefaultShortDescription = "Gestion de tareas y sesion (MySQL)"
	DefaultLongDescription  = "Tareas con persistencia en MySQL, login y sesion desde la terminal."
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
