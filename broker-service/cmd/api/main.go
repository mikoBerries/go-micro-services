package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type config struct{}

func main() {
	app := config{}
	// make server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("starting broker services ad :%s\n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
