package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)


func CreateDB() *sql.DB {
	dsn := "file:mailing-list.db?cache=shared"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(fmt.Sprintf("Could not get a database connection: %v\n", err))
	}

	return db
}


func CreateSchema(db *sql.DB) {
	stmt := `
		CREATE TABLE IF NOT EXISTS emails (
			id INTEGER NOT NULL PRIMARY KEY,
			email TEXT NOT NULL,
			date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
			date_verified DATETIME DEFAULT NULL,
			subscribed BOOLEAN DEFAULT true
		);
	`
	fmt.Printf("Creating emails table:\n%v\n", stmt)
	_, err := db.Exec(stmt)

	if err != nil {
		panic(fmt.Sprintf("Could not create emails table: %v\n", err))
	}
}