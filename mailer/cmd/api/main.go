package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const connection_port = "80"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port", connection_port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", connection_port),
		Handler: app.routes(),
	}

	possible_error := server.ListenAndServe()
	if possible_error != nil {
		log.Panic(possible_error)
	}

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mail_obj := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return mail_obj
}
