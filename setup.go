package main

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

const (
	configPath = "./config.json"
)

type configStruct struct {
	Port               int
	SlackClientID      string
	SlackClientSecret  string
	SlackSigningSecret string
}

var config configStruct

func setup() error {
	f, err := os.Open(configPath)
	if err != nil {
		return errors.Wrap(err, "open config file")
	}

	err = json.NewDecoder(f).Decode(&config)
	defer f.Close()
	if err != nil {
		return errors.Wrap(err, "decode config contents")
	}

	return nil
}
