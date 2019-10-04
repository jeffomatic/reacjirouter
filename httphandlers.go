package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/tokenstore"
)

func extractFormParam(r *http.Request, key string) (string, bool) {
	param, ok := r.Form[key]
	if !ok || len(param) != 1 {
		return "", false
	}
	return param[0], true
}

// TODO: validate state param
func handleSlackOauth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	code, ok := extractFormParam(r, "code")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := slack.NewClient("").GetAccessToken(slack.AccessTokenArgs{
		ClientID:     config.SlackClientID,
		ClientSecret: config.SlackClientSecret,
		Code:         code,
	})
	if err != nil {
		log.Println("slack oauth.access error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.TeamID == "" || resp.AccessToken == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenstore.Add(resp.TeamID, resp.AccessToken)

	log.Println("received access token for team", resp.TeamID, resp.AccessToken)
	w.Write([]byte("Access token received for team " + resp.TeamID))
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

func handleSlashCommand(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	teamID, ok := extractFormParam(r, "team_id")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	text, ok := extractFormParam(r, "text")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message := processCommand(teamID, text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slack.SlashCommandResponse{
		ResponseType: "ephemeral",
		Text:         message,
	})
}
