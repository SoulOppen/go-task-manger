package auth

import (
	"bufio"
	"fmt"
	"os"
)

func SignUp() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Binevenido al sistema de conexion")
	fmt.Println("¿Cual es tu nombre de usuario?")
	scanner.Scan()
	username := scanner.Text()
	_ = username
	fmt.Println("¿Cual es tu clave?")
	scanner.Scan()
	password := scanner.Text()
	_ = password
	fmt.Println("Fuiste registrado con exito")
}
