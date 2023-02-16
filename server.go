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
	clients []*Client
	broadcast chan string
}

func (c *Client) handleConnection() {
	defer c.conn.Close()
	c.username = c.readUsername()
	fmt.Println("Client", c.username, "connected.")

	var roomNames []string
	for roomName := range rooms {
		roomNames = append(roomNames, roomName)
	}
	fmt.Fprintf(c.conn, "Available rooms: %s\n", strings.Join(roomNames, ", "))

	fmt.Fprintf(c.conn, "Enter the name of the room you want to join: ")
    roomName, _ := bufio.NewReader(c.conn).ReadString('\n')
    roomName = strings.TrimSpace(roomName)
    
    // Join the specified room
    room, ok := rooms[roomName]
    if !ok {
        fmt.Fprintf(c.conn, "Invalid room name. Goodbye!\n")
        return
    }
    room.join(c)

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
	for _, room := range rooms {
		if containsClient(room.clients, origin.conn) {
			for _, c := range room.clients {
				if c.conn != origin.conn {
					c.conn.Write([]byte(message))
				}
			}
		}
	}

	for _, client := range clients {
		if client == origin {
			continue
		}
		client.conn.Write([]byte(message))
	}
}

func createRooms() {
	roomList := []string{"General", "Programming", "Gaming", "Chess", "Music", "Misc", "The Ratway", "File transfer"}
	for _, roomName := range roomList {
		room := &Room {
			name: roomName,
			clients: make([]*Client, 0),
			broadcast: make(chan string),
		}
		rooms[roomName] = room
	}
}

func (r *Room) join(c *Client) {
    r.clients = append(r.clients, c)
    fmt.Printf("Client %s joined room %s\n", c.username, r.name)
}

func (r *Room) leave(c *Client) {
    for i, client := range r.clients {
        if client == c {
            r.clients = append(r.clients[:i], r.clients[i+1:]...)
            fmt.Printf("Client %s left room %s\n", c.username, r.name)
            break
        }
    }
}

func containsClient(clients []*Client, conn net.Conn ) bool {
	for _, client := range clients {
		if client.conn == conn {
			return true
		}
	}
	return false
}

var clients []*Client

var rooms = make(map[string]*Room)

func main() {
	createRooms()

	l, err := net.Listen("tcp", "0.0.0.0:1491")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer l.Close()

	fmt.Println("Waiting for Daenerys to finish saying all her titles..")

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
