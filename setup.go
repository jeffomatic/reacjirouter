package main

import "github.com/jeffomatic/reacjirouter/slack"

const (
	slackAPIToken = "changeme"
)

var (
	appUserID string
	teamURL   string
)

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
