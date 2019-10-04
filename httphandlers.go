package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/tokenstore"
)

const oauthTemplate = `
<html>
  <head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/3.0.1/github-markdown.min.css">
    <style>
      .markdown-body {
        box-sizing: border-box;
        min-width: 200px;
        max-width: 980px;
        margin: 0 auto;
        padding: 45px;
      }

      @media (max-width: 767px) {
        .markdown-body {
          padding: 15px;
        }
      }
    </style>
  </head>
  <body>
    <article class="markdown-body">
      <h1>Reacji Router installed!</h1>
      <pre lang="no-highlight"><code>{{respJSON}}</code></pre>
    </article>
  </body>
</html>
`

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

	var resp slack.OauthAccessResponse
	err = slack.NewClient("").Call(slack.OauthAccess, slack.OauthAccessRequest{
		AuthorizationType: "grant",
		ClientID:          config.SlackClientID,
		ClientSecret:      config.SlackClientSecret,
		Code:              code,
	}, &resp)
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

	respJSON, _ := json.MarshalIndent(resp, "", "  ") // should never fail
	out := strings.Replace(oauthTemplate, "{{respJSON}}", strings.TrimSpace(string(respJSON)), 1)
	w.Write([]byte(out))
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

	if e.T != "event_callback" {
		log.Printf("/slack/event: received unknown payload type %q", e.T)
		return
	}

	if e.Event.T != "reaction_added" {
		log.Printf("/slack/event: received unknown event type %q", e.Event.T)
		return
	}

	c := newTeamClient(e.TeamID)
	if c == nil {
		return
	}

	err = handleReactionAdded(c, e.Event.Reaction, e.Event.Item.ChannelID, e.Event.Item.Timestamp)
	if err != nil {
		log.Printf("/slack/event: error handling reaction event: %s", err)
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
