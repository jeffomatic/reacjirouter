package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/store"
)

var (
	spaceSplitter *regexp.Regexp
)

func init() {
	spaceSplitter = regexp.MustCompile(`\s+`)
}

func handleIM(c *teamClient, channelID, userID, text string) error {
	text = strings.TrimSpace(text)
	tokens := spaceSplitter.Split(text, 3)
	var err error

	switch strings.ToLower(tokens[0]) {
	case "add":
		return handleAddCommand(c, channelID, userID, tokens)

	case "list":
		return handleListCommand(c, channelID, userID, tokens)

	case "help":
		return handleHelpCommand(c, channelID, userID, tokens)

	default:
		err = c.sendEphemeralMessage(userID, channelID, "Sorry, I didn't recognize that command! Type \"help\" for instructions.")
		if err != nil {
			return err
		}
	}

	return nil
}

func handleReactionAdded(c *teamClient, emoji string, channelID string, timestamp string) error {
	targetChannel, ok := store.Get(c.teamID, emoji)
	if !ok {
		return nil
	}

	teamURL, err := c.teamURL()
	if err != nil {
		return err
	}

	text := slack.MessageLink{teamURL, channelID, timestamp}.String()
	err = c.sendMessage(targetChannel, text)
	if err != nil {
		return err
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
