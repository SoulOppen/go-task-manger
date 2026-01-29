package name

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Name() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("No se pudo cargar env")
		os.Exit(1)
	}
	name := os.Getenv("NAME")
	name = strings.ReplaceAll(name, " ", "-")
	return name
}
