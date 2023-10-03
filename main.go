package main

import (
	"fmt"
	"log"
	"net"
	"theratway/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {	
	db, err := server.InitDB("mydatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	server.CreateRooms()

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
		client := server.NewClient(conn)
		server.Clients = append(server.Clients, client)
		go client.HandleConnection(db)
	}
}