package server

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"theratway/casino"
	"theratway/main"
	"theratway/room"
)

var clientRoles = make(map[string][]*Client)
var rolesList = []string{"Programmer", "Gamer", "Rat", "Gopher", "Paranormal Investigator","Mod"}
var Clients []*Client

func isValidRole(role string, roles []string) bool {
    for _, r := range roles {
        if r == role {
            return true
        }
    }
    return false
}

type Client struct {
	Conn     net.Conn
	Username string
	Password string
	Role     string
	Player   *casino.Player
}

func NewClient(conn net.Conn) *Client {
	return &Client{Conn: conn}
}

func (c *Client) listenForCommand(message string, db *sql.DB) bool {
	if !strings.HasPrefix(message, "/") {
		return false
	}

	command := strings.Split(message[1:], " ")
	switch command[0] {
	case "join":
		roomName := strings.TrimSpace(command[1])
		room, ok := room.rooms[roomName]
		if !ok {
			room := room.makeRoom(roomName)
			room.rooms[roomName] = room
		}
		room.Join(c)
		joinmsg := "%s has joined"
		room.broadcastMessage(fmt.Sprintf(joinmsg, c.Username), c)
		return false
	case "leave":
		{
			if room, ok := room.clientRooms[c]; ok {
				room.Leave(c)
				ext := "%s has left"
				room.broadcastMessage(fmt.Sprintf(ext, c.Username), c)

				c.showRooms()
				c.Conn.Write([]byte("\nWelcome to The Ratway.\n"))
			}
			return false
		}
	case "role":
		if len(command) < 2 {
			c.Conn.Write([]byte("Please specify a role.\n"))
			return false
		}
		role := command[1]
		if !room.isValidRole(role, room.rolesList) {
			c.Conn.Write([]byte("Invalid role.\n"))
			return false
		}
		room.clientRoles[role] = append(room.clientRoles[role], c)
		c.Conn.Write([]byte(fmt.Sprintf("You have been assigned the role of %s.\n", role)))
		// add user's role to database
		err := c.updateRole(db, role, c.Username)
		if err != nil {
			log.Fatal(err)
		}
		return false
	case "help":
		c.Conn.Write([]byte("List of commands:\n/help for list of commands\n/join <room name> to join a room.\n/leave to leave a room.\n/list to list all users in a room\n/role <role> to assign yourself a role.\n"))
		return false
	case "list":
		roomName := room.clientRooms[c].name
		room := room.getRoom(roomName)
		listUsers(c, room)
	default:
		// When the command doesn't exist
		fmt.Println("Unknown command:", command[0])
		return true
	}
	return true
}

func (c *Client) HandleConnection(db *sql.DB) {
	defer c.Conn.Close()

	login(db, c)

	for {
		message, err := bufio.NewReader(c.Conn).ReadString('\x00')
		if err != nil {
			fmt.Println("Client", c.Username, "disconnected.")
			return
		}
		if message == "" {
			continue
		}

		listenFor := c.listenForCommand(filterString(message), db)
		_ = listenFor

		if strings.HasPrefix(message, "/") {
			continue
		}

		if room, ok := room.clientRooms[c]; ok {
			message = c.Username + ": " + message
			fmt.Println(message)
			room.broadcastMessage(message, c)
		}
	}
}

func (c *Client) showRooms() {
	c.Conn.Write([]byte("Available rooms:\n"))

	for _, name := range room.static_rooms {
		c.Conn.Write([]byte("-" + fmt.Sprint(name) + "\n"))
	}
}

func (c *Client) readUsername() string {
	c.Conn.Write([]byte("Enter username: "))
	username, _ := bufio.NewReader(c.Conn).ReadString('\x00')
	username = strings.TrimSpace(username)
	if username == "" {
		username = "anon"
	}
	return username
}

func (c *Client) readPassword() string {
	c.Conn.Write([]byte("Enter password: "))
	password, _ := bufio.NewReader(c.Conn).ReadString('\x00')
	password = strings.TrimSpace(password)
	return password
}

func listUsers(c *Client, r *room.Room) {
	userList := ""
	for _, user := range r.Members {
		userList += user + "\n"
	}
	c.Conn.Write([]byte(userList))
}

func CreateRooms() {
	for _, roomName := range room.static_rooms {
		room := &room.Room {
			Name : roomName,
			clients: []Client{},
			Members: make(map[net.Conn]string),
			broadcast: make(chan string),
		}
		room.rooms[roomName] = room
	}
}