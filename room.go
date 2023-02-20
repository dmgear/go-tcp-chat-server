package main

import (
	"fmt"
	"net"
)

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