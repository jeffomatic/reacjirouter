package main

import "github.com/jeffomatic/reacjirouter/slack"

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
