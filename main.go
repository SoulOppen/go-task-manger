package main

import (
	"fmt"
)

type User struct {
	id   string
	name string
	coin string
}

func (u *User) setUser(id, name, coin string) {
	u.id = id
	u.name = name
	u.coin = coin
}

var user User

func main() {
	userName := welcome()
	user.setUser("129192", userName, "ajsknajsdjk")
	fmt.Println("Welcome,", userName)
	fmt.Println(user)
}
