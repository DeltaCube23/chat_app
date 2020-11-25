package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

//Server is
type Server struct {
	Password          string
	ClientConnections map[string]net.Conn
	mu                sync.RWMutex
}

//NewServer instance
func NewServer(password string) *Server {
	return &Server{
		Password:          password,
		ClientConnections: make(map[string]net.Conn),
	}
}

//ListenForMessages waits to hear from clients
func (s *Server) ListenForMessages(conn net.Conn, name string) {
	for {
		select {
		default:
			buf := make([]byte, 256)
			learn := bufio.NewReader(conn)

			length := 0
			for length < 256 {
				alpha, err := learn.ReadByte()
				if err != nil {
					return
				}
				buf[length] = alpha
				// $ is the delimiter
				if alpha == '$' {
					break // to quit loop
				}
				length++
			}
			msg := string(buf[:length])
			//parse the string
			msg = strings.Trim(msg, "\r\n")
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])

			switch cmd {
			case "/pm":
				s.pm(conn, name, args)
			case "/broad":
				s.broad(conn, name, args)
			case "/quit":
				s.quit(conn, name)
			}
		}
	}
}

//for private message
func (s *Server) pm(conn net.Conn, name string, args []string) {
	msg := strings.Join(args[2:], " ")
	s.mu.RLock()
	to := s.ClientConnections[args[1]]
	s.mu.RUnlock()
	msg = "private from " + name + " : " + msg + " $"
	to.Write([]byte(msg))
}

//for broadcast
func (s *Server) broad(conn net.Conn, name string, args []string) {
	msg := strings.Join(args[1:], " ")
	s.mu.RLock()
	for user, id := range s.ClientConnections {
		if user != name {
			msg = "general from " + name + " : " + msg + " $"
			id.Write([]byte(msg))
		}
	}
	s.mu.RUnlock()
}

//to quit
func (s *Server) quit(conn net.Conn, name string) {
	log.Printf("client %s has left the chat", name)
	s.mu.Lock()
	to := s.ClientConnections[name]
	delete(s.ClientConnections, name) //erase detailts from the map
	s.mu.Unlock()
	to.Write([]byte("Good Bye$"))
	to.Close()
}

//ManageClient for server
func (s *Server) ManageClient(conn net.Conn) {
	buf := make([]byte, 256)
	learn := bufio.NewReader(conn)

	length := 0
	for length < 256 {
		alpha, err := learn.ReadByte()
		if err != nil {
			return
		}
		buf[length] = alpha
		// $ is the delimiter
		if alpha == '$' {
			break // to quit loop
		}
		length++
	}
	//to parse the string
	msg := string(buf[:length])
	msg = strings.Trim(msg, "\r\n")
	args := strings.Split(msg, " ")
	cmd := strings.TrimSpace(args[0])
	var reply string

	if cmd == "/auth" {
		if args[1] == s.Password { // check for correct password
			s.mu.RLock()
			_, ok := s.ClientConnections[args[2]] // check if username exists
			s.mu.RUnlock()
			if ok == true {
				reply = "Username Already taken. Password was correct $"
				conn.Write([]byte(reply))
			} else {
				log.Printf("new client has joined: %s", args[2])
				reply = "you are now an authenticated user $"
				conn.Write([]byte(reply))
				s.mu.Lock()
				s.ClientConnections[args[2]] = conn
				s.mu.Unlock()
				go s.ListenForMessages(conn, args[2])
			}
		} else {
			reply = "Wrong Password $"
			conn.Write([]byte(reply))
		}
	}

}

//ListenForConnections from client
func (s *Server) ListenForConnections(joined chan net.Conn, listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("failed to accept connection: %s", err.Error())
		return
	}
	joined <- conn
}

//Run starts the server
func (s *Server) Run(ctx context.Context) {
	listener, err := net.Listen("tcp", ":"+"8888")
	if err != nil {
		fmt.Printf("unable to start server: %s", err.Error())
		return
	}
	defer listener.Close()

	joined := make(chan net.Conn)

	for {
		go s.ListenForConnections(joined, listener)
		select {
		case <-ctx.Done(): // when server exits
			s.blackout()
			time.Sleep(5 * time.Second)
			fmt.Println("Terminating Server...")
			return
		case conn := <-joined: // when new client is connected
			go s.ManageClient(conn)
		}
	}
}

//terminate all clients before server exits
func (s *Server) blackout() {
	s.mu.Lock()
	for name, id := range s.ClientConnections {
		delete(s.ClientConnections, name)
		id.Write([]byte("Good Bye$"))
		fmt.Printf("Terminating Client %s...\n", name)
	}
	s.mu.Unlock()
}
