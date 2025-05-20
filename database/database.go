package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"solar-garlic-mailing-list/model"

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
			email TEXT NOT NULL UNIQUE,
			date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
			date_verified DATETIME DEFAULT NULL,
			subscribed BOOLEAN DEFAULT true
		);
	`
	slog.Info("Creating emails table", "stmt", stmt)
	_, err := db.Exec(stmt)

	if err != nil {
		slog.Error("Error creating emails table", "err", err)
		panic(fmt.Sprintf("Could not create emails table: %v\n", err))
	}
}

func CreateEmail(db *sql.DB, email string) (*model.Email, error) {
	var createdEmail = &model.Email{}
	stmt := `
		INSERT INTO emails (email) VALUES ($1)
		RETURNING id, email, date_created, date_verified, subscribed
	`
	params := []any{email}

	slog.Info("Inserting email into database", "stmt", stmt, "params", params)
	err := db.QueryRow(stmt, params...).Scan(
		&createdEmail.ID,
		&createdEmail.Email,
		&createdEmail.DateCreated,
		&createdEmail.DateVerified,
		&createdEmail.Subscribed,
	)

	if err != nil {
		slog.Error("Error inserting email into database", "stmt", stmt, "params", params)
		return nil, err
	}

	return createdEmail, nil
}

func ListEmails(db *sql.DB, email string) ([]model.Email, error) {
	var emails []model.Email

	stmt, params := getListEmailsQuery(email)

	slog.Info("Selecting emails from database", "stmt", stmt, "params", params)

	rows, err := db.Query(stmt, params...)
	if err != nil {
		slog.Error("Error selecting email list from database", "stmt", stmt, "params", params)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var email model.Email

		scanErr := rows.Scan(
			&email.ID,
			&email.Email,
			&email.DateCreated,
			&email.DateVerified,
			&email.Subscribed,
		)

		if scanErr != nil {
			slog.Error("Error scanning sql rows", "stmt", stmt, "params", params, "err", scanErr)
			return nil, err
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func getListEmailsQuery(email string) (string, []any) {
	if email == "" {
		stmt := `
			SELECT  
				id,
				email,
				date_created,
				date_verified,
				subscribed	
			FROM emails;	
		`
		params := []any{}
		return stmt, params
	}

	stmt := `
		SELECT	
			id,
			email,
			date_created,
			date_verified,
			subscribed
		FROM emails
		where email = $1;
	`
	params := []any{email}
	return stmt, params
}
