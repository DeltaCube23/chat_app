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
	Conn     net.Conn
}

//NewClient instance is created
func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
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
			fmt.Printf("\n%s", msg)
		}
	}
}

func (c *Client) getClientMessage(recv chan string) {
	msg, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return
	}
	recv <- msg
}

//HandleClient reads messages from client
func (c *Client) HandleClient() {
	recv := make(chan string)
	for {
		go c.getClientMessage(recv)
		select {
		case msg := <-recv:
			//parse the string
			//msg = strings.Trim(msg, "\r\n")
			//args := strings.Split(msg, " ")
			//cmd := strings.TrimSpace(args[0])
			//send respective command details to server channel

			c.Conn.Write([]byte(msg))
		}
	}
}

//Run the client side
func (c *Client) Run() {
	conn, err := net.Dial("tcp", ":"+"8080")
	if err != nil {
		fmt.Printf("Unable to connect to server : %s", err.Error())
	}

	c.Conn = conn
	msg := "/auth " + c.Password + " " + c.Username + "$"
	c.Conn.Write([]byte(msg))

	go c.HandleClient()
	go c.HandleServer()
	select {}
}

func (c *Client) err(err error) {
	c.Conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *Client) msg(msg string) {
	c.Conn.Write([]byte("> " + msg + "\n"))
}
