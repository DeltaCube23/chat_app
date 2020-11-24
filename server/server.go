package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
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

			msg = strings.Trim(msg, "\r\n")
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])

			switch cmd {
			case "/pm":
				s.pm(conn, args)
			case "/broad":
				s.broad(conn, args)
			case "/quit":
				s.quit(conn, args)
			}
		}
	}
}

func (s *Server) pm(conn net.Conn, args []string) {
	msg := strings.Join(args[3:], " ")
	s.mu.RLock()
	to := s.ClientConnections[args[2]]
	s.mu.RUnlock()
	msg = "private from " + args[1] + " : " + msg + " $"
	to.Write([]byte(msg))
}

func (s *Server) broad(conn net.Conn, args []string) {
	msg := strings.Join(args[2:], " ")
	s.mu.RLock()
	for name, id := range s.ClientConnections {
		if args[1] != name {
			msg = "general from " + args[1] + " : " + msg + " $"
			id.Write([]byte(msg))
		}
	}
	s.mu.RUnlock()
}

func (s *Server) quit(conn net.Conn, args []string) {
	log.Printf("client %s has left the chat", args[1])
	s.mu.Lock()
	to := s.ClientConnections[args[1]]
	delete(s.ClientConnections, args[1]) //erase detailts from the map
	s.mu.Unlock()
	to.Write([]byte("Good Bye $"))
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
	msg := string(buf[:length])

	msg = strings.Trim(msg, "\r\n")
	args := strings.Split(msg, " ")
	cmd := strings.TrimSpace(args[0])
	var reply string

	if cmd == "/auth" {
		if args[1] == s.Password {
			s.mu.RLock()
			_, ok := s.ClientConnections[args[2]]
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
	for {
		select {
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("failed to accept connection: %s", err.Error())
				continue
			}
			joined <- conn
		}
	}
}

//Run starts the server
func (s *Server) Run() {
	listener, err := net.Listen("tcp", ":"+"8080")
	if err != nil {
		fmt.Printf("unable to start server: %s", err.Error())
		return
	}
	defer listener.Close()

	joined := make(chan net.Conn)
	go s.ListenForConnections(joined, listener)

	for {
		select {
		case conn := <-joined:
			go s.ManageClient(conn)
		}
	}
}
