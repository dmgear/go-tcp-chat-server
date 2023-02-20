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
	role string
}

type Room struct {
	clients []Client
	name string
	Members map[net.Conn]string
	broadcast chan string
}

func (r *Room) Join(client *Client) {
	r.Members[client.conn] = client.username
	clientRooms[client] = r
}

func (r * Room) Leave(client *Client) {
	delete(r.Members, client.conn)
	delete(clientRooms, client)
}

func (r * Room) broadcastMessage(message string, origin *Client) {
	for conn := range r.Members {
		if conn == origin.conn {
			continue
		}
		conn.Write([]byte(fmt.Sprintf("%s", message)))
	}
}

func (c *Client) giveRole(role string, client *Client) {
	c.role = role
	clientRoles[role] = append(clientRoles[role], client)
}

func (c *Client) listenForCommand(message string) {
	if !strings.HasPrefix(message, "/") {
		return
	}

	command := strings.Split(message[1:], " ")
	switch command[0] {
	case "join":
		roomName := strings.TrimSpace(strings.Join(command[1:], " "))
		room, ok := rooms[roomName]
		if !ok {
			room = &Room{
				name:      roomName,
				Members:   make(map[net.Conn]string),
				broadcast: make(chan string),
			}
			rooms[roomName] = room
		}
		room.Join(c)
		joinmsg := "%s has joined"
		room.broadcastMessage(fmt.Sprintf(joinmsg, c.username), c)
	case "leave":
		if room, ok := clientRooms[c]; ok {
			room.Leave(c)
			ext := "%s has left"
			room.broadcastMessage(fmt.Sprintf(ext, c.username), c)
			c.showRooms()
		}
	default:
		// When the command doesn't exist
		fmt.Println("Unknown command:", command[0])
	}
}

func (c *Client) showRooms() {
	c.conn.Write([]byte("Available rooms:\n"))

	for _, name := range static_rooms { 
		c.conn.Write([]byte("-" + fmt.Sprint(name) + "\n"))
	}
}

var static_rooms = []string{"#General", "#Programming", "#Gaming", "#Music", "#Misc", "#The Ratway", "#File transfer"}

func (c *Client) handleConnection() {
	defer c.conn.Close()
	c.username = strings.TrimSpace(c.readUsername())
	c.username = c.username[:len(c.username)-1]
	fmt.Println("Client", c.username, "connected.")

	c.showRooms()

	for {
		message, err := bufio.NewReader(c.conn).ReadString('\x00')
		if err != nil {
			fmt.Println("Client", c.username, "disconnected.")
			return
		}
		if message == "" {
			continue
		}

	c.listenForCommand(message)
	if room, ok := clientRooms[c]; ok {
		message = c.username + ": " + message
		fmt.Println(message)
		room.broadcastMessage(message, c)
	}
	}
}

func (c *Client) readUsername() string {
    c.conn.Write([]byte("Enter username: "))
    username, _ := bufio.NewReader(c.conn).ReadString('\x00')
    username = strings.TrimSpace(username)
    if username == "" {
        username = "anon"
    }
    return username
}

func createRooms() {
	for _, roomName := range static_rooms {
		room := &Room {
			name : roomName,
			clients: []Client{},
			Members: make(map[net.Conn]string),
			broadcast: make(chan string),
		}
		rooms[roomName] = room
	}
}
var clients []*Client

var rooms = make(map[string]*Room)
var clientRooms = make(map[*Client]*Room)
var clientRoles = make(map[string][]*Client)

func createRolesMap() {
	var roles = []string{"Programmer", "Gamer", "Rat", "Gopher", "Mod"}
	for _, role := range roles {
		clientRoles[role] = []*Client{}
	}
}

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