package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type slackEvent struct {
	T        string `json:"type"`
	UserID   string `json:"user"`
	Reaction string
	Item     struct {
		T         string `json:"type"`
		ChannelID string `json:"channel"`
		Timestamp string `json:"ts"`
	}
}

type slackEventPayload struct {
	T         string `json:"type"`
	Event     slackEvent
	Challenge string // for URL verification
}

func handleSlackEvent(w http.ResponseWriter, r *http.Request) {
	v, ok := r.Header["Content-Type"]
	if !ok || v[0] != "application/json" {
		fmt.Println(r.Header)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var e slackEventPayload
	err := decoder.Decode(&e)
	if err != nil {
		log.Println("/slack/event: error decoding JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch e.T {
	case "url_verification":
		fmt.Fprintf(w, e.Challenge)

	case "event_callback":
		switch e.Event.T {
		case "reaction_added":
			log.Printf(
				"TODO: reaction %q user %q channel %q timestamp %q",
				e.Event.Reaction,
				e.Event.UserID,
				e.Event.Item.ChannelID,
				e.Event.Item.Timestamp,
			)
		default:
			log.Printf("/slack/event: received unknown event type %q", e.Event.T)
		}

	default:
		log.Printf("/slack/event: received unknown payload type %q", e.T)
	}
}

func main() {
	router := mux.NewRouter()
	router.Path("/slack/event").Methods("POST").HandlerFunc(handleSlackEvent)

	handler := handlers.LoggingHandler(os.Stdout, router)

	port := 1234
	fmt.Println("Starting server on port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Fatal(err)
	}
}
