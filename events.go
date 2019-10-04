package main

import (
	"github.com/jeffomatic/reacjirouter/routestore"
	"github.com/jeffomatic/reacjirouter/slack"
)

func handleReactionAdded(c *teamClient, emoji string, channelID string, timestamp string) error {
	targetChannel, ok := routestore.Get(c.teamID, emoji)
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
