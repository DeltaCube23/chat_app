package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DeltaCube23/chat_app/client"
	"github.com/DeltaCube23/chat_app/server"
)

var (
	password string
	username string
	address  string
)

func gracexit(cancel context.CancelFunc) {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM) //look for termination signals
	sig := <-sigchan
	log.Printf("shutdown : %v", sig)
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	if os.Args[1] == "s" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)
		fmt.Println("Enter port number to connect to : ")
		fmt.Scanln(&address)

		go gracexit(cancel)
		s := server.NewServer(password, address)
		s.Run(ctx) // Start Server
	} else if os.Args[1] == "c" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)
		fmt.Println("Enter UserName : ")
		fmt.Scanln(&username)
		fmt.Println("Enter port number to connect to : ")
		fmt.Scanln(&address)

		c := client.NewClient(username, password, address)
		c.Run() // Start Client
	}
}
