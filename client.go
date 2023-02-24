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
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) listenForCommand(message string) bool {
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
		return false
	default:
		// When the command doesn't exist
		fmt.Println("Unknown command:", command[0])
		return true
	}
}

func (c *Client) handleConnection(db *sql.DB) {
	defer c.conn.Close()

	// get username and password from client
	c.username = strings.TrimSpace(c.readUsername())
	c.username = filterString(c.username)
	
	// check if user already exists in the database
    userExists, err := checkUserExists(db, c.username)
    if err != nil {
        log.Fatal(err)
        return
    }
    if userExists {
        fmt.Println("User", c.username, "already exists.")
        return
    } 
	if !userExists {
	// prompt user to create new account with email
	c.conn.Write([]byte("No existing account with that username, please sign up "))
	c.username = strings.TrimSpace(c.readUsername())
	c.username = filterString(c.username)
	c.password = strings.TrimSpace(c.readPassword())
	c.password = filterString(c.password)
	fmt.Println("Client", c.username, "connected.")

	c.showRooms() // display list of rooms to user
	}
    // add new user to the database
    err = addUser(db, c.username, c.password)
    if err != nil {
        log.Fatal(err)
        return
    }

	for {
		message, err := bufio.NewReader(c.conn).ReadString('\x00')
		if err != nil {
			fmt.Println("Client", c.username, "disconnected.")
			return
		}
		if message == "" {
			continue
		}

		listenFor := c.listenForCommand(filterString(message))
		_ = listenFor
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