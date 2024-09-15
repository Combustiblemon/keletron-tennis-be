package database

import (
	"database/sql"
	"log"
)

var tables = []string{
	`CREATE TABLE IF NOT EXISTS "reservations" (
	"_id"	TEXT NOT NULL UNIQUE,
	"type"	TEXT NOT NULL,
	"datetime"	TEXT NOT NULL,
	"duration"	INTEGER NOT NULL,
	"owner"	TEXT NOT NULL,
	"court"	TEXT NOT NULL,
	"status"	TEXT DEFAULT 'PENDING',
	"paid"	TEXT DEFAULT 'FALSE',
	"notes"	TEXT,
	FOREIGN KEY("owner") REFERENCES "users"("_id"),
  FOREIGN KEY("court") REFERENCES "courts"("_id"),
	PRIMARY KEY("_id")
);`,

	`CREATE TABLE IF NOT EXISTS "users" (
	"_id"	TEXT NOT NULL UNIQUE,
	"role"	TEXT NOT NULL DEFAULT 'USER',
	"email"	TEXT NOT NULL,
	"password"	TEXT,
	"session"	TEXT UNIQUE,
	"accountType"	TEXT DEFAULT 'PASSWORD',
	"FCMTokens"	TEXT,
	"resetKey"	TEXT,
	PRIMARY KEY("_id")
);`,

	`CREATE TABLE IF NOT EXISTS "courts" (
	"_id" TEXT NOT NULL UNIQUE,
	"name"	TEXT NOT NULL,
	"type"	TEXT NOT NULL,
	"reservationStartTime"	TEXT NOT NULL DEFAULT '09:00',
	"reservationEndTime"	TEXT NOT NULL DEFAULT '21:00',
	"reservationDuration"	INTEGER NOT NULL DEFAULT 90,
	PRIMARY KEY("_id")
);`,

	`CREATE TABLE IF NOT EXISTS "announcements" (
	"_id" TEXT NOT NULL UNIQUE,
	"body"	TEXT,
	"title"	TEXT NOT NULL,
	"validUntil"	TEXT NOT NULL,
	"visible"	TEXT,
	PRIMARY KEY("_id")
);`,

	`CREATE TABLE IF NOT EXISTS "court_reserved_times" (
	"_id" TEXT NOT NULL UNIQUE,
	"court"	TEXT NOT NULL,
	"duration"	INTEGER NOT NULL,
	"type"	TEXT NOT NULL,
	"repeat"	TEXT NOT NULL,
	"days"	TEXT NOT NULL,
	"notes"	TEXT NOT NULL,
	"daysNotApplied"	TEXT NOT NULL DEFAULT "[]",
  FOREIGN KEY("court") REFERENCES "courts"("_id"),
	PRIMARY KEY("_id")
);`,
}

func setupTables(db *sql.DB) error {
	for _, statement := range tables {
		_, err := db.Exec(statement)

		if err != nil {
			return err
		}
	}

	return nil
}

var db *sql.DB

func Setup() error {
	log.Println("Setting up database")

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	err = setupTables(db)

	if err != nil {
		return err
	}

	log.Println("Database setup complete")

	return nil
}

func Teardown() error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
