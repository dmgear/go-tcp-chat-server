package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	data     chan []byte
	username string
	outgoing chan string
	role     string
}

func (c *Client) giveRole(role string, client *Client) {
	c.role = role
	clientRoles[role] = append(clientRoles[role], client)
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

func (c *Client) handleConnection() {
	defer c.conn.Close()
	c.username = strings.TrimSpace(c.readUsername())
	c.username = filterString(c.username)
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

		listenFor := c.listenForCommand(filterString(message))
		_ = listenFor
		if room, ok := clientRooms[c]; ok {
			message = c.username + ": " + message
			fmt.Println(message)
			if listenFor {
				room.broadcastMessage(message, c)
			}
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