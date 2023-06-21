package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const webPort = "80"

type Config struct {
	Mailer Mail
}

func main() {
	log.Println("Processing mail-services")
	time.Sleep(15 * time.Second)
	log.Println("Getting all env config")
	app := Config{
		Mailer: createMail(),
	}
	log.Printf("%+v \n", app.Mailer)
	// make mail service
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("starting mail-services at :%s\n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// createMail for creating Mail and setting from ENV to SMTP server
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	// all seting for SMTP server
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("EMAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}
	return m
}
