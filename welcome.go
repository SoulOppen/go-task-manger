package main

import (
	"fmt"
)

func welcome() string {
	var username string

	fmt.Println("Welcome to Task Manager API.")
	fmt.Println("Your tasks, organized and under control.")
	fmt.Println("What is your username?")
	fmt.Scanln(&username)

	return username
}
