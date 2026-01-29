package auth

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SignUp() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Binevenido al sistema de conexion")
	fmt.Println("¿Cual es tu nombre de usuario?")
	scanner.Scan()
	username := scanner.Text()
	fmt.Println("¿Cual es tu clave?")
	scanner.Scan()
	password := scanner.Text()
	_ = password
	confDir, err := os.UserConfigDir()
	if err != nil {
		os.Exit(1)
	}
	dir := filepath.Join(confDir, "usuarios.txt")
	data, err := os.ReadFile(dir)
	if err != nil {
		os.Exit(1)
	}
	contenido := strings.Split(string(data), "\n")
	contenido = append(contenido, username)
	texto := strings.Join(contenido, "\n")
	err = os.WriteFile(dir, []byte(texto), 0644)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Fuiste registrado con exito")
}
