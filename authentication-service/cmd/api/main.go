package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MikoBerries/go-micro-services/authentication-service/data"

	// db driver
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	Repo   data.Repository //populate db and models
	Client *http.Client    // for mocking roundTrip
}

func main() {
	log.Println("Starting authentication-service")
	// conn to db postgress
	conn := connecToDB()
	if conn == nil {
		log.Panic("Cant connect to Postgres")
	}

	app := Config{
		Client: &http.Client{},
	}
	// make serv
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("authentication-service listen and serve (%s) \n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Ping connection db
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connecToDB() *sql.DB {
	var counts int64
	interval := 5 * time.Second
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready")
			counts++
		} else {
			log.Println("Connected to Postgres: " + dsn)
			return connection
		}
		if counts > 10 {
			log.Panic(err)
			return nil
		}
		log.Printf("Retrying connection in %v second, number of retry - %v", interval, counts)
		time.Sleep(interval)
		continue
	}
	// return nil
}
