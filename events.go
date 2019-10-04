package main

import (
	"regexp"
	"strings"

	"github.com/jeffomatic/reacjirouter/routestore"
	"github.com/jeffomatic/reacjirouter/slack"
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

	switch strings.ToLower(tokens[0]) {
	case "add":
		return handleAddCommand(c, channelID, userID, tokens)

	case "list":
		return handleListCommand(c, channelID, userID, tokens)

	case "help":
		return handleHelpCommand(c, channelID, userID, tokens)

	default:
		err := c.sendEphemeralMessage(userID, channelID, "Sorry, I didn't recognize that command! Type \"help\" for instructions.")
		if err != nil {
			return err
		}
	}

	return nil
}

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
