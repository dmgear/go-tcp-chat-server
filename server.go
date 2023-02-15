package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
	data chan []byte
	username string
	outgoing chan string
}

func (c *Client) handleConnection() {
	defer c.conn.Close()
	c.username = c.readUsername()
	fmt.Println("Client", c.username, "connected.")

	for {
		message, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			fmt.Println("Client", c.username, "disconnected.")
			return
		}
		if message == "" {
			continue
		}
		message = c.username + ": " + message
		broadcastMessage(message, c)
	}
}

func (c *Client) readUsername() string {
	c.conn.Write([]byte("enter username: "))
	username, _ := bufio.NewReader(c.conn).ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "" {
		username = "anon"
	}
	return username
}

func broadcastMessage(message string, origin *Client) {
	for _, client := range clients {
		if client == origin {
			continue
		}
		client.conn.Write([]byte(message))
	}
}

var clients []*Client

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:<port>")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer l.Close()

	fmt.Println("Waiting for Danaerys to finish saying all her titles..")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		client := &Client{conn: conn}
		clients = append(clients, client)
		go client.handleConnection()
	}
}
