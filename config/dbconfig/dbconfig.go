package dbconfig

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB(databaseURL string) *sql.DB {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database : %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed : %v", err)
	}

	log.Println("Connected to database")

	return db
}
