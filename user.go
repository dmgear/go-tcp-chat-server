package main

import "database/sql"


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