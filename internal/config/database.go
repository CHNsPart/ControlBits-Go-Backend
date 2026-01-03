package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dsn := "host=localhost port=5432 user=controlbits password=controlbits123 dbname=controlbits_db sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect database:", err)
	}

	fmt.Println("✅ Database connected")
	return db
}
