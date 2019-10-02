package main

import (
	"encoding/json"
	"os"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/pkg/errors"
)

const (
	configPath = "./config.json"
)

var (
	slackAPIToken string
	appUserID     string
	teamURL       string
)

func setup() error {
	f, err := os.Open(configPath)
	if err != nil {
		return errors.Wrap(err, "open config file")
	}

	var config struct{ SlackAPIToken string }
	err = json.NewDecoder(f).Decode(&config)
	defer f.Close()
	if err != nil {
		return errors.Wrap(err, "decode config contents")
	}

	slackAPIToken = config.SlackAPIToken

	var resp slack.AuthTestResponse
	err = slack.NewClient(slackAPIToken).Call("auth.test", nil, &resp)
	if err != nil {
		return err
	}

	appUserID = resp.UserID
	teamURL = resp.URL

	return nil
}
