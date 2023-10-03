package room

import (
	"net"
	"theratway/main"
	"theratway/server"
	"strings"
)

var static_rooms = []string{"#General", "#Programming", "#Gaming", "#Music", "#Paranormal", "#Misc", "#The Ratway", "#File transfer"}
var rooms = make(map[string]*Room)
var clientRooms = make(map[*server.Client]*Room)

type Room struct {
	clients []server.Client
	Name string
	Members map[net.Conn]string
	broadcast chan string
}

func (r *Room) Join(client *server.Client) {
	r.Members[client.Conn] = client.Username
	 clientRooms[client] = r
}

func (r * Room) Leave(client *server.Client) {
	delete(r.Members, client.Conn)
	delete(clientRooms, client)
}

func (r * Room) broadcastMessage(message string, origin *server.Client) {
	for conn := range r.Members {
		if conn == origin.Conn {
			continue
		}
		conn.Write([]byte(message))
	}
}

func getRoom(name string) *Room {
	for _, room := range rooms {
		if name == name {
			return room
		}
	}
	return nil
}

func makeRoom(roomName string) Room {
	room := &Room{
		Name:      roomName,
		Members:   make(map[net.Conn]string),
		broadcast: make(chan string),
	}
	return *room
}


