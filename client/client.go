package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

//Client is
type Client struct {
	Username string
	Password string
	Address  string
	Conn     net.Conn
}

//NewClient instance is created
func NewClient(username, password, address string) *Client {
	return &Client{
		Username: username,
		Password: password,
		Address:  address,
	}
}

func (c *Client) getServerMessage(recv chan string) {
	buf := make([]byte, 256)
	learn := bufio.NewReader(c.Conn)

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
	recv <- msg
}

//HandleServer listens for messages from server
func (c *Client) HandleServer() {
	recv := make(chan string)
	for {
		go c.getServerMessage(recv)
		select {
		case msg := <-recv:
			fmt.Printf("\n>> %s\n>> ", msg)
			if msg == "Good Bye" {
				os.Exit(0)
			}
		}
	}
}

//Read message of client from STDIN
func (c *Client) getClientMessage(recv chan string) {
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	msg = msg + "$"
	recv <- msg
}

//HandleClient reads messages from client
func (c *Client) HandleClient() {
	recv := make(chan string)
	for {
		go c.getClientMessage(recv)
		select {
		case msg := <-recv:
			c.Conn.Write([]byte(msg))
		}
		fmt.Printf(">> ")
	}
}

//Run the client side
func (c *Client) Run() {
	service := ":" + c.Address
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Printf("unable to resolve tcpaddress: %s", err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Printf("Unable to connect to server : %s", err.Error())
	}

	c.Conn = conn
	msg := "/auth " + c.Password + " " + c.Username + "$"
	c.Conn.Write([]byte(msg)) //To Authenticate the User

	go c.HandleClient()
	go c.HandleServer()
	select {}
}
