package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	slackAPIURLPrefix = "https://slack.com/api/"
	slackAPIToken     = "changeme"
)

var emojiToChannel map[string]string

func init() {
	emojiToChannel = make(map[string]string)

	// TEMP: hardcode some mappings
	emojiToChannel["grin"] = "#general"
}

type slackEvent struct {
	ID       string `json:"event_id"`
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
	TeamID    string `json:"team_id"`
	Event     slackEvent
	Challenge string // for URL verification
}

func handleReacji(reacji string, srcChannel string, timestamp string) error {
	channel, ok := emojiToChannel[reacji]
	if !ok {
		return nil
	}

	reqBody, err := json.Marshal(struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
		AsUser  bool   `json:"as_user"`
	}{
		channel,
		fmt.Sprintf("reacji: %s channel: %s ts: %s", reacji, srcChannel, timestamp), // TODO
		true,
	})
	if err != nil {
		panic("error encoding JSON: " + err.Error()) // should never happen
	}

	req, err := http.NewRequest(http.MethodPost, slackAPIURLPrefix+"chat.postMessage", bytes.NewReader(reqBody))
	if err != nil {
		panic("error building request: " + err.Error()) // should never happen
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Authorization"] = []string{"Bearer " + slackAPIToken}

	c := new(http.Client)
	resp, err := c.Do(req)
	if err != nil {
		return errors.Wrap(err, "chat.postMessage transport error")
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("chat.postMessage bad status: %s", resp.Status)
	}

	respBody := struct {
		Ok    bool
		Error string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return errors.Wrap(err, "chat.postMessage response body decode error")
	}

	if !respBody.Ok {
		return fmt.Errorf("chat.postMessage error: %s", respBody.Error)
	}

	return nil
}

func handleSlackEvent(w http.ResponseWriter, r *http.Request) {
	v, ok := r.Header["Content-Type"]
	if !ok || v[0] != "application/json" {
		fmt.Println(r.Header)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var e slackEventPayload
	err := json.NewDecoder(r.Body).Decode(&e)
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
			err = handleReacji(e.Event.Reaction, e.Event.Item.ChannelID, e.Event.Item.Timestamp)
			if err != nil {
				log.Printf("/slack/event: error handling reacji: %s", err)
			}
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
