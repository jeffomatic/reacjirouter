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

func slackAPICall(method string, body interface{}) error {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "error encoding JSON")
	}

	req, err := http.NewRequest(
		http.MethodPost,
		slackAPIURLPrefix+method,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return errors.Wrap(err, "error preparing request")
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Authorization"] = []string{"Bearer " + slackAPIToken}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return errors.Wrap(err, "transport error")
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	respBody := struct {
		Ok    bool
		Error string
	}{}
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return errors.Wrap(err, "response body decode error")
	}

	if !respBody.Ok {
		return fmt.Errorf("Slack API error: %s", respBody.Error)
	}

	return nil
}

func handleReacji(reacji string, srcChannel string, timestamp string) error {
	channel, ok := emojiToChannel[reacji]
	if !ok {
		return nil
	}

	err := slackAPICall("chat.postMessage", struct {
		Channel string `json:"channel"`
		Text    string `json:"text"`
		AsUser  bool   `json:"as_user"`
	}{
		channel,
		fmt.Sprintf("reacji: %s channel: %s ts: %s", reacji, srcChannel, timestamp), // TODO
		true,
	})
	if err != nil {
		return errors.Wrap(err, "chat.postMessage error")
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
