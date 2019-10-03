package main

import (
	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/tokens"
)

type teamClient struct {
	*slack.Client
	teamID string

	// don't access the following directly
	authTestCache *slack.AuthTestResponse
}

func newTeamClient(teamID string) *teamClient {
	token, ok := tokens.Get(teamID)
	if !ok {
		return nil
	}

	return &teamClient{Client: slack.NewClient(token), teamID: teamID}
}

func (c *teamClient) teamURL() (string, error) {
	if c.authTestCache != nil {
		return c.authTestCache.URL, nil
	}

	var resp slack.AuthTestResponse
	err := c.Client.Call("auth.test", nil, &resp)
	if err != nil {
		return "", err
	}

	c.authTestCache = &resp
	return c.authTestCache.URL, nil
}

func (c *teamClient) userID() (string, error) {
	if c.authTestCache != nil {
		return c.authTestCache.UserID, nil
	}

	var resp slack.AuthTestResponse
	err := c.Client.Call("auth.test", nil, &resp)
	if err != nil {
		return "", err
	}

	c.authTestCache = &resp
	return c.authTestCache.UserID, nil
}

func (c *teamClient) sendMessage(channelID, text string) error {
	return c.Client.Call(
		"chat.postMessage",
		slack.ChatPostMessageRequest{channelID, text, true},
		nil,
	)
}

func (c *teamClient) sendEphemeralMessage(userID, channelID, text string) error {
	return c.Client.Call(
		"chat.postEphemeral",
		slack.ChatPostEphemeralRequest{userID, channelID, text, true},
		nil,
	)
}
