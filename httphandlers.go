package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeffomatic/reacjirouter/slack"
)

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

	if e.T == "url_verification" {
		fmt.Fprintf(w, e.Challenge)
		return
	}

	// All other events require a valid token
	c := newTeamClient(e.TeamID)
	if c == nil {
		return
	}

	switch e.T {
	case "url_verification":
		fmt.Fprintf(w, e.Challenge)

	case "event_callback":
		switch e.Event.T {
		case "message":
			switch e.Event.ChannelType {
			case "im":
				appUserID, err := c.userID()
				if err != nil {
					log.Println("error fetching user ID:", err)
					break
				}

				// Don't respond to messages the app itself generates
				if e.Event.UserID == appUserID {
					break
				}

				err = handleIM(c, e.Event.ChannelID, e.Event.UserID, e.Event.Text)
				if err != nil {
					log.Printf("/slack/event: error handling IM event: %s", err)
				}

			default:
				log.Printf("/slack/event: can't handle message channel type: %s", e.Event.ChannelType)
			}

		case "reaction_added":
			err = handleReactionAdded(c, e.Event.Reaction, e.Event.Item.ChannelID, e.Event.Item.Timestamp)
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
