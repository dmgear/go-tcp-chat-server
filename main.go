package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

func isValidRole(role string, roles []string) bool {
    for _, r := range roles {
        if r == role {
            return true
        }
    }
    return false
}

func filterString(message string) string {
	message = message[:len(message)-1]
	return message
}

var static_rooms = []string{"#General", "#Programming", "#Gaming", "#Music", "#Misc", "#The Ratway", "#File transfer"}

var clients []*Client

var rooms = make(map[string]*Room)

var clientRooms = make(map[*Client]*Room)

var clientRoles = make(map[string][]*Client)

var rolesList = []string{"Programmer", "Gamer", "Rat", "Gopher", "Mod"}

func main() {	
	db, err := InitDB("mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
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