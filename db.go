package main

import (
	"database/sql"
)



func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./mydatabase.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS USERS (
		ID INTEGER PRIMARY KEY,
		USERNAME TEXT,
		PASSWORD TEXT,
		EMAIL TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
	return db, nil
}
