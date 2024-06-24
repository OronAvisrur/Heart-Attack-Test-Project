package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const conn_port = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authenication")

	// Connecting to database
	connection := connectToDB()

	if connection == nil {
		log.Panic("Can't connect to database!")
	}

	// Setting up config
	app := Config{
		DB:     connection,
		Models: data.New(connection),
	}

	// Setting up server connection
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", conn_port),
		Handler: app.routes(),
	}

	//------------
	possible_error := server.ListenAndServe()
	if possible_error != nil {
		log.Panic(possible_error)
	}

}

// This function responsible to open connection to given database
func openDB(dsn string) (*sql.DB, error) {
	// Trying to open connect
	db, possible_error := sql.Open("pgx", dsn)
	if possible_error != nil {
		return nil, possible_error
	}

	// Trying to send ping
	possible_error = db.Ping()
	if possible_error != nil {
		return nil, possible_error
	}

	// No errors return the database
	return db, nil
}

// This function responsible to the connection
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, possible_error := openDB(dsn)
		if possible_error != nil {
			log.Println("Can't connect to database")
			count++
		} else {
			log.Println("Connected successfully to database")
			return connection
		}

		// If we still can't connect after 10 tries stop tyring
		if count > 10 {
			log.Println(possible_error)
			return nil
		}

		// Waiting 2 second before trying to connect again
		time.Sleep(2 * time.Second)
		continue
	}
}
