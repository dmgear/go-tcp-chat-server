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

type Room struct {
	name string
	Members map[net.Conn]string
	broadcast chan string
}

func (c *Client) handleConnection() {
	defer c.conn.Close()
	c.username = c.readUsername()
	fmt.Println("Client", c.username, "connected.")

	for {
		message, err := bufio.NewReader(c.conn).ReadString('\x00')
		if err != nil {
			fmt.Println("Client", c.username, "disconnected.")
			return
		}
		if message == "" {
			continue
		}
		message = c.username + ": " + message
		fmt.Println(message)
		broadcastMessage(message, c)
	}
}

func (c *Client) readUsername() string {
	c.conn.Write([]byte("enter username: "))
	username, _ := bufio.NewReader(c.conn).ReadString('\x00')
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

func createRooms() {
	//roomList := []string{"General", "Programming", "Gaming", "Music", "Misc", "The Ratway", "File transfer"}
}

var clients []*Client

var rooms = make(map[string]*Room)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:1491")

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