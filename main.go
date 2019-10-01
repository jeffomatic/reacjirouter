package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/store"
)

const (
	slackAPIURLPrefix = "https://slack.com/api/"
	slackAPIToken     = "changeme"
)

var (
	spaceSplitter *regexp.Regexp
	appUserID     string
	teamURL       string
)

func init() {
	spaceSplitter = regexp.MustCompile(`\s+`)
}

func sendMessage(channelID, text string) error {
	return slack.NewClient(slackAPIToken).Call(
		"chat.postMessage",
		slack.ChatPostMessageRequest{channelID, text, true},
		nil,
	)
}

func sendEphemeralMessage(userID, channelID, text string) error {
	return slack.NewClient(slackAPIToken).Call(
		"chat.postEphemeral",
		slack.ChatPostEphemeralRequest{userID, channelID, text, true},
		nil,
	)
}

func handleIM(teamID, channelID, userID, text string) error {
	text = strings.TrimSpace(text)
	tokens := spaceSplitter.Split(text, 3)
	var err error

	switch strings.ToLower(tokens[0]) {
	case "add":
		if len(tokens) != 3 {
			err = sendEphemeralMessage(userID, channelID, `Could not understand add command`)
			if err != nil {
				return err
			}

			return nil // this is a user error
		}

		emoji, ok := slack.ExtractEmoji(tokens[1])
		if !ok {
			log.Println("emoji error", tokens[1])
			err = sendEphemeralMessage(userID, channelID, `Could not understand add command`)
			if err != nil {
				return err
			}

			return nil // this is a user error
		}

		targetChannelID, ok := slack.ExtractChannelID(tokens[2])
		if !ok {
			log.Println("channel error", tokens[2])
			err = sendEphemeralMessage(userID, channelID, `Could not understand add command`)
			if err != nil {
				return err
			}

			return nil // this is a user error
		}

		// TODO: handle the following error conditions
		// - bot not in channel

		store.Add(teamID, emoji, targetChannelID)

		text := fmt.Sprintf("Okay, I'll send all messages with the :%s: reacji to <#%s>.", emoji, targetChannelID)
		err = sendMessage(channelID, text)
		if err != nil {
			return err
		}

	case "help":
		err = sendEphemeralMessage(userID, channelID, `
*Instructions*

Add a new reaction route
> add :emoji: #channel

Show help
> help
`)
		if err != nil {
			return err
		}

	default:
		err = sendEphemeralMessage(userID, channelID, "Sorry, I didn't recognize that command! Type \"help\" for instructions.")
		if err != nil {
			return err
		}
	}

	return nil
}

func handleReactionAdded(teamID string, emoji string, channelID string, timestamp string) error {
	targetChannel, ok := store.Get(teamID, emoji)
	if !ok {
		return nil
	}

	text := slack.MessageLink{teamURL, channelID, timestamp}.String()
	err := sendMessage(targetChannel, text)
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

	switch e.T {
	case "url_verification":
		fmt.Fprintf(w, e.Challenge)

	case "event_callback":
		switch e.Event.T {
		case "message":
			switch e.Event.ChannelType {
			case "im":
				// Don't respond to messages the app itself generates
				if e.Event.UserID == appUserID {
					break
				}

				err = handleIM(e.TeamID, e.Event.ChannelID, e.Event.UserID, e.Event.Text)
				if err != nil {
					log.Printf("/slack/event: error handling IM event: %s", err)
				}

			default:
				log.Printf("/slack/event: can't handle message channel type: %s", e.Event.ChannelType)
			}

		case "reaction_added":
			err = handleReactionAdded(e.TeamID, e.Event.Reaction, e.Event.Item.ChannelID, e.Event.Item.Timestamp)
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

func setupTeamData() error {
	var resp slack.AuthTestResponse
	err := slack.NewClient(slackAPIToken).Call("auth.test", nil, &resp)
	if err != nil {
		return err
	}

	appUserID = resp.UserID
	teamURL = resp.URL

	return nil
}

func main() {
	err := setupTeamData()
	if err != nil {
		log.Fatalln("could not fetch app user ID, error:", err)
	}

	router := mux.NewRouter()
	router.Path("/slack/event").Methods("POST").HandlerFunc(handleSlackEvent)

	handler := handlers.LoggingHandler(os.Stdout, router)

	port := 1234
	log.Println("Starting server on port", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Fatal(err)
	}
}
