package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DeltaCube23/chat_app/client"
	"github.com/DeltaCube23/chat_app/server"
)

var (
	password string
	username string
)

func gracexit(cancel context.CancelFunc) {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM) //look for termination signals
	sig := <-sigchan
	log.Printf("shutdown : %v", sig)
	time.Sleep(5 * time.Second)
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	if os.Args[1] == "s" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)

		go gracexit(cancel)
		s := server.NewServer(password)
		s.Run(ctx) // Start Server
	} else if os.Args[1] == "c" {
		fmt.Println("Enter Server Password : ")
		fmt.Scanln(&password)
		fmt.Println("Enter UserName : ")
		fmt.Scanln(&username)

		c := client.NewClient(username, password)
		c.Run() // Start Client
	}
}
