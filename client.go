package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	username string
	password string
	role     string
	hand     []Card
	pile     []Card
	points   int
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) listenForCommand(message string, db *sql.DB) bool {
	if !strings.HasPrefix(message, "/") {
		return false
	}

	command := strings.Split(message[1:], " ")
	switch command[0] {
	case "join":
		roomName := strings.TrimSpace(command[1])
		roomName = filterString(roomName)
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
		return false
	case "leave":
		{
			if room, ok := clientRooms[c]; ok {
				room.Leave(c)
				ext := "%s has left"
				room.broadcastMessage(fmt.Sprintf(ext, c.username), c)

				c.showRooms()
				c.conn.Write([]byte("\nWelcome to the lobby.\n"))

			}
			return false
		}
	case "role":
		if len(command) < 2 {
			c.conn.Write([]byte("Please specify a role.\n"))
			return false
		}
		role := command[1]
		if !isValidRole(role, rolesList) {
			c.conn.Write([]byte("Invalid role.\n"))
			return false
		}
		clientRoles[role] = append(clientRoles[role], c)
		c.conn.Write([]byte(fmt.Sprintf("You have been assigned the role of %s.\n", role)))
		// add user's role to database
		err := c.updateRole(db, role, c.username)
		if err != nil {
			log.Fatal(err)
		}
		return false
	case "help":
		c.conn.Write([]byte("List of commands:\n/help for list of commands\n/join <room name> to join a room.\n/leave to leave a room.\n/list to list all users in a room\n/role <role> to assign yourself a role.\n"))
		return false
	case "list":
		roomName := clientRooms[c].name
		room := getRoom(roomName)
		listUsers(c, room)
	default:
		// When the command doesn't exist
		fmt.Println("Unknown command:", command[0])
		return true
	}
	return true
}

func (c *Client) handleConnection(db *sql.DB) {
	defer c.conn.Close()

	login(db, c)

	for {
		message, err := bufio.NewReader(c.conn).ReadString('\x00')
		if err != nil {
			fmt.Println("Client", c.username, "disconnected.")
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

		if room, ok := clientRooms[c]; ok {
			message = c.username + ": " + message
			fmt.Println(message)
			room.broadcastMessage(message, c)
		}
	}
}

func (c *Client) showRooms() {
	c.conn.Write([]byte("Available rooms:\n"))

	for _, name := range static_rooms {
		c.conn.Write([]byte("-" + fmt.Sprint(name) + "\n"))
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

func (c *Client) readPassword() string {
	c.conn.Write([]byte("Enter password: "))
	password, _ := bufio.NewReader(c.conn).ReadString('\x00')
	password = strings.TrimSpace(password)
	return password
}

func listUsers(c *Client, r *Room) {
	userList := ""
	for _, user := range r.Members {
		userList += user + "\n"
	}
	c.conn.Write([]byte(userList))
}

func getRoom(name string) *Room {
	for _, room := range rooms {
		if room.name == name {
			return room
		}
	}
	return nil
}
