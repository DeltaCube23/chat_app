package main

import (
	"fmt"
	"os"

	"github.com/DeltaCube23/chat_app/client"
	"github.com/DeltaCube23/chat_app/server"
)

var (
	password string
	username string
)

func main() {
	if os.Args[1] == "s" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)

		s := server.NewServer(password)
		s.Run()
	} else if os.Args[1] == "c" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)
		fmt.Println("Enter UserName : ")
		fmt.Scanln(&username)
		c := client.NewClient(username, password)
		c.Run()
	}
}
