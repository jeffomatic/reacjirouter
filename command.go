package main

import (
	"fmt"
	"strings"

	"github.com/jeffomatic/reacjirouter/routes"
	"github.com/jeffomatic/reacjirouter/slack"
)

func handleAddCommand(c *teamClient, channelID, userID string, tokens []string) error {
	if len(tokens) != 3 {
		err := c.sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	emoji, ok := slack.ExtractEmoji(tokens[1])
	if !ok {
		err := c.sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	targetChannelID, ok := slack.ExtractChannelID(tokens[2])
	if !ok {
		err := c.sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	// TODO: handle the following error conditions
	// - bot not in channel

	routes.Add(c.teamID, emoji, targetChannelID)

	text := fmt.Sprintf("Okay, I'll send all messages with the :%s: reacji to <#%s>.", emoji, targetChannelID)
	return c.sendMessage(channelID, text)
}

func handleListCommand(c *teamClient, channelID, userID string, tokens []string) error {
	if len(tokens) > 1 {
		err := c.sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}
	}

	text := `No reacji configured.`
	list := routes.List(c.teamID)
	if len(list) > 0 {
		lines := make([]string, len(list))
		for _, pair := range list {
			lines = append(lines, fmt.Sprintf(":%s: :point_right: <#%s>", pair.Emoji, pair.ChannelID))
		}
		text = strings.Join(lines, "\n")
	}

	return c.sendMessage(channelID, text)
}

func handleHelpCommand(c *teamClient, channelID, userID string, tokens []string) error {
	text := `
*Instructions*

Add a new reaction route
> add :emoji: #channel

List reaction routes on this team
> list

Show help
> help
`
	return c.sendEphemeralMessage(userID, channelID, text)
}
