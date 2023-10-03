package server

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
    "theratway/main"
)

func filterString(message string) string {
	message = message[:len(message)-1]
	return message
}

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
    query := "INSERT INTO USERS(username, password) VALUES(?, ?)"

    // Execute the query to insert username and password into db
    _, err := db.Exec(query, username, password)
    if err != nil {
        return err
    }
    return nil
}

func checkPassword(db *sql.DB, username string, password string) (bool, error) {
    var dbPassword string
    err := db.QueryRow("SELECT PASSWORD FROM USERS WHERE USERNAME=?", username).Scan(&dbPassword)
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
	c.Username = strings.TrimSpace(c.readUsername())
	c.Username = filterString(c.Username)
	
	// check if user already exists in the database
    userExists, err := checkUserExists(db, c.Username)
    if err != nil {
        log.Fatal(err)
        return
    }
    if userExists {
		c.Conn.Write([]byte("User already exists, please log in"))
        fmt.Println("User", c.Username, "already exists, please log in.")
		c.Password = strings.TrimSpace(c.readPassword())
		c.Password = filterString(c.Password)
		result, err := checkPassword(db, c.Username, c.Password)
		if err != nil {
			log.Fatal(err)
		}

		for !result {
			c.Conn.Write([]byte("incorrect password"))
			c.Password = strings.TrimSpace(c.readPassword())
			c.Password = filterString(c.Password)
			result, err = checkPassword(db, c.Username, c.Password)
			if err != nil {
				log.Fatal(err)
			}
		}
		c.showRooms()
    } 

	if !userExists {
	// prompt user to create new account with email
	c.Conn.Write([]byte("No existing account with that username, please sign up "))
	c.Password = strings.TrimSpace(c.readPassword())
	c.Password = filterString(c.Password)
	fmt.Println("Client", c.Username, "connected.")

	// add new user to the database if they dont already exist
    err = addUser(db, c.Username, c.Password)
    if err != nil {
        log.Fatal(err)
        return
    }
	// display list of rooms to user
	c.showRooms() 
	}
}