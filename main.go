package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatalln("could not load config, error:", err)
	}

	router := mux.NewRouter()
	router.Path("/slack/command").Methods("POST").HandlerFunc(handleSlashCommand)
	router.Path("/slack/event").Methods("POST").HandlerFunc(handleSlackEvent)
	router.Path("/slack/oauth").Methods("GET").HandlerFunc(handleSlackOauth)

	handler := handlers.LoggingHandler(os.Stdout, router)

	log.Println("Starting server on port", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler)
	if err != nil {
		log.Fatal(err)
	}
}
