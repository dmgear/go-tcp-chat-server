package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
    conn net.Conn
    data chan []byte
    username string
    outoging chan string
}

func (c *Client) handleConnection() {
    defer c.conn.Close()
    c.username = c.readUsername()
    fmt.Println("Client", c.username, "connected.")

    c.conn.Write([]byte("Available rooms:\x00"))
    for name := range rooms {
        c.conn.Write([]byte("- " + name + "\x00"))
    }

    reader := bufio.NewReader(c.conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client", c.username, "disconnected.")
            } else {
                fmt.Println("Read error:", err)
            }
            return
        }
        message = strings.TrimSpace(message)
        if message == "" {
            continue
        }
        message = c.username + ": " + message
        fmt.Println(message)
        broadcastMessage(message, c)
    }
}

func (c *Client) readUsername() string {
    username, _ := bufio.NewReader(c.conn).ReadString('\n')
    username = strings.TrimSpace(username)
    if username == "" {
        username = "anon"
    }
    return username
}

func broadcastMessage(message string) {
    for _, client := range clients {
        client.conn.Write([]byte(message))
    }
}

var clients []*Client

func main() {
    l, err := net.Listen("tcp", "0.0.0.0:1491")

    if err != nil {
        log.Fatal(err)
        return
    }

    defer l.Close()

    fmt.Println("Waiting for Danaerys to finish saying all her titles..")

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