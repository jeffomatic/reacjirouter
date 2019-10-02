package main

import (
	"fmt"
	"strings"

	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/store"
)

func handleAddCommand(teamID, channelID, userID string, tokens []string) error {
	if len(tokens) != 3 {
		err := sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	emoji, ok := slack.ExtractEmoji(tokens[1])
	if !ok {
		err := sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	targetChannelID, ok := slack.ExtractChannelID(tokens[2])
	if !ok {
		err := sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}

		return nil // this is a user error
	}

	// TODO: handle the following error conditions
	// - bot not in channel

	store.Add(teamID, emoji, targetChannelID)

	text := fmt.Sprintf("Okay, I'll send all messages with the :%s: reacji to <#%s>.", emoji, targetChannelID)
	return sendMessage(channelID, text)
}

func handleListCommand(teamID, channelID, userID string, tokens []string) error {
	if len(tokens) > 1 {
		err := sendEphemeralMessage(userID, channelID, `Could not understand add command`)
		if err != nil {
			return err
		}
	}

	text := `No reacji configured.`
	list := store.List(teamID)
	if len(list) > 0 {
		lines := make([]string, len(list))
		for _, pair := range list {
			lines = append(lines, fmt.Sprintf(":%s: :point_right: <#%s>", pair.Emoji, pair.ChannelID))
		}
		text = strings.Join(lines, "\n")
	}

	return sendMessage(channelID, text)
}

func handleHelpCommand(channelID, userID string, tokens []string) error {
	text := `
*Instructions*

Add a new reaction route
> add :emoji: #channel

List reaction routes on this team
> list

Show help
> help
`
	return sendEphemeralMessage(userID, channelID, text)
}
