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

func (c *teamClient) ensureAuthTestCache() error {
	if c.authTestCache != nil {
		return nil
	}

	resp := new(slack.AuthTestResponse)
	err := c.Client.Call(slack.AuthTest, nil, resp)
	if err != nil {
		return err
	}

	c.authTestCache = resp
	return nil
}

func (c *teamClient) teamURL() (string, error) {
	err := c.ensureAuthTestCache()
	if err != nil {
		return "", err
	}

	return c.authTestCache.URL, nil
}

func (c *teamClient) userID() (string, error) {
	err := c.ensureAuthTestCache()
	if err != nil {
		return "", err
	}

	return c.authTestCache.UserID, nil
}
