package main

import (
	"database/sql"
	"log"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS USERS (
		USERNAME TEXT,
		PASSWORD TEXT,
		ROLE TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
	return db, nil
}

func (c *Client) updateRole(db *sql.DB, role string, username string) error {
	c.role = role
	c.username = username
	query := "UPDATE USERS SET ROLE=? WHERE USERNAME=?"
	
	_, err := db.Exec(query, role, username)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
