package auth

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func usuariosFile() string {
	confDir, err := os.UserConfigDir()
	if err != nil {
		os.Exit(1)
	}
	return filepath.Join(confDir, "usuarios.txt")
}

func leerUsuarios() []string {
	dir := usuariosFile()

	data, err := os.ReadFile(dir)
	if err != nil {
		return []string{}
	}

	return strings.FieldsFunc(string(data), func(r rune) bool {
		return r == '\n' || r == '\r'
	})
}

func guardarUsuarios(usuarios []string) {
	dir := usuariosFile()

	texto := strings.Join(usuarios, "\n")
	err := os.WriteFile(dir, []byte(texto), 0644)
	if err != nil {
		os.Exit(1)
	}
}

func SignUp() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Bienvenido al sistema de conexión")
	fmt.Println("¿Cuál es tu nombre de usuario?")
	scanner.Scan()
	username := scanner.Text()

	fmt.Println("¿Cuál es tu clave?")
	scanner.Scan()
	password := scanner.Text()
	_ = password

	usuarios := leerUsuarios()

	if slices.Contains(usuarios, username) {
		fmt.Println("El usuario ya existe")
		os.Exit(1)
	}

	usuarios = append(usuarios, username)
	guardarUsuarios(usuarios)

	fmt.Println("Fuiste registrado con éxito")
}

func Login() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("¿Cuál es tu nombre de usuario?")
	scanner.Scan()
	username := scanner.Text()

	usuarios := leerUsuarios()

	if !slices.Contains(usuarios, username) {
		fmt.Println("No existe usuario")
		os.Exit(1)
	}

	fmt.Println("Login exitoso")
}
