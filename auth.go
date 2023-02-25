package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)


func checkUserExists(db *sql.DB, username string) (bool, error) {
    // Prepare a SELECT query to check if the user exists
    query := "SELECT COUNT(*) FROM users WHERE username = ?"
	
    // Execute the query with the given parameters
    var count int
    err := db.QueryRow(query, username).Scan(&count)
    if err != nil {
        return false, err
    }
    // If count is greater than zero, the user exists
    return count > 0, nil
}

func addUser(db *sql.DB, username string, password string) error {
    // Prepare an INSERT query to add a new user
    query := "INSERT INTO users(username, password) VALUES(?, ?)"

    // Execute the query with the given parameters
    _, err := db.Exec(query, username, password)
    if err != nil {
        return err
    }
    return nil
}

func checkPassword(db *sql.DB, username string, password string) (bool, error) {
    var dbPassword string
    err := db.QueryRow("SELECT PASSWORD FROM USERS WHERE USERNAME = ?", username).Scan(&dbPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return false, nil
        }
        return false, err
    }
    return dbPassword == password, nil
}

func login(db *sql.DB, c *Client) {
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
		c.conn.Write([]byte("User already exists, please log in"))
        fmt.Println("User", c.username, "already exists, please log in.")
		c.password = strings.TrimSpace(c.readPassword())
		c.password = filterString(c.password)
		result, err := checkPassword(db, c.username, c.password)
		if err != nil {
			log.Fatal(err)
		}
		for !result {
			c.conn.Write([]byte("incorrect password"))
			c.password = strings.TrimSpace(c.readPassword())
			c.password = filterString(c.password)
			result, err = checkPassword(db, c.username, c.password)
			if err != nil {
				log.Fatal(err)
			}
		}
		c.showRooms()
    } 
	
	if !userExists {
	// prompt user to create new account with email
	c.conn.Write([]byte("No existing account with that username, please sign up "))
	c.username = strings.TrimSpace(c.readUsername())
	c.username = filterString(c.username)
	c.password = strings.TrimSpace(c.readPassword())
	c.password = filterString(c.password)
	fmt.Println("Client", c.username, "connected.")

	// add new user to the database if they dont already exist
    err = addUser(db, c.username, c.password)
    if err != nil {
        log.Fatal(err)
        return
    }
	// display list of rooms to user
	c.showRooms() 
	}
}