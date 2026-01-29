package name

import (
	"os"

	"github.com/joho/godotenv"
)

func Name() (string, error) {
	err := godotenv.Load()
	if err == nil {
		return "", err
	}
	name := os.Getenv("NAME")
	return name, nil
}
