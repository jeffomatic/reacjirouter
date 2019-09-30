package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jeffomatic/reacjirouter/slack"
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

func buildMessageLink(teamID string, channelID string, timestamp string) (string, error) {
	var resp slack.TeamInfoResponse

	// TODO: cache this
	err := slack.NewClient(slackAPIToken).Call("team.info", nil, &resp)
	if err != nil {
		return "", errors.Wrap(err, "team domain fetch")
	}

	// Format:
	// https://{domain}.slack.com/archives/{channel ID}/p{timestamp with period stripped}
	return fmt.Sprintf(
		"https://%s.slack.com/archives/%s/p%s",
		resp.Team.Domain,
		channelID,
		strings.Replace(timestamp, ".", "", 1),
	), nil
}

func handleReactionAdded(emoji string, teamID string, channelID string, timestamp string) error {
	channel, ok := emojiToChannel[emoji]
	if !ok {
		return nil
	}

	message, err := buildMessageLink(teamID, channelID, timestamp)
	if err != nil {
		return errors.Wrap(err, "building message link")
	}

	err = slack.NewClient(slackAPIToken).Call(
		"chat.postMessage",
		slack.ChatPostMessageRequest{channel, message, true},
		nil,
	)
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

	var e slack.EventPayload
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
			err = handleReactionAdded(e.Event.Reaction, e.TeamID, e.Event.Item.ChannelID, e.Event.Item.Timestamp)
			if err != nil {
				log.Printf("/slack/event: error handling reaction event: %s", err)
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
