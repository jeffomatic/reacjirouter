package main

import (
	"log"
	"net/http"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/tokens"
)

// TODO: validate state param
func handleSlackOauth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	codeParam, ok := r.Form["code"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	if len(codeParam) != 1 {
		w.WriteHeader(http.StatusBadRequest)
	}

	resp, err := slack.NewClient("").GetAccessToken(slack.AccessTokenArgs{
		ClientID:     config.SlackClientID,
		ClientSecret: config.SlackClientSecret,
		Code:         codeParam[0],
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

	tokens.Add(resp.TeamID, resp.AccessToken)

	log.Println("received access token for team", resp.TeamID, resp.AccessToken)
	w.Write([]byte("Access token received for team " + resp.TeamID))
}
