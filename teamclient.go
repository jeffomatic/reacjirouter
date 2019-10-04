package main

import (
	"github.com/jeffomatic/reacjirouter/slack"
	"github.com/jeffomatic/reacjirouter/tokenstore"
)

type teamClient struct {
	*slack.Client
	teamID string

	// don't access the following directly
	authTestCache *slack.AuthTestResponse
}

func newTeamClient(teamID string) *teamClient {
	token, ok := tokenstore.Get(teamID)
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
