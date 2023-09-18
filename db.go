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

	createUserTable := `CREATE TABLE IF NOT EXISTS USERS (
		USERID INTEGER PRIMARY KEY,
		USERNAME TEXT,
		PASSWORD TEXT
	);`
	_, err = db.Exec(createUserTable)
	if err != nil {
		panic(err)
	}

	createRolesTable := `CREATE TABLE IF NOT EXISTS ROLES (
		ROLEID INTEGER PRIMARY KEY,
		ROLENAME TEXT,
		USERID INTEGER,
		FOREIGN KEY (USERID) REFERENCES USERS(USERID)
	);`

	_, err = db.Exec(createRolesTable)
	if err != nil {
		panic(err)
	}
	return db, nil
}

func (c *Client) updateRole(db *sql.DB, role string, username string) error {
	c.role = role
	c.username = username
	query := "INSERT INTO ROLES (ROLENAME, USERID) VALUES (?, (SELECT USERID FROM USERS WHERE USERNAME=?))"
	
	_, err := db.Exec(query, role, username)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
