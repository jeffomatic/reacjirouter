package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const defaultURLPrefix = "https://slack.com/api"

type Client struct {
	URLPrefix   string
	AccessToken string
}

func NewClient(token string) *Client {
	return &Client{URLPrefix: defaultURLPrefix, AccessToken: token}
}

func (c *Client) Call(method string, body interface{}, respBody interface{}) error {
	reqBody := []byte(`{}`)
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "json encode request body")
		}
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.URLPrefix+"/"+method,
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return errors.Wrap(err, "HTTP request creation")
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Authorization"] = []string{"Bearer " + c.AccessToken}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return errors.Wrap(err, "executing HTTP request")
	}
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "reading HTTP response body")
	}

	errCheck := struct {
		Ok    bool
		Error string
	}{}
	err = json.Unmarshal(respBytes, &errCheck)
	if err != nil {
		return errors.Wrap(err, "response body error-check decode")
	}
	if !errCheck.Ok {
		return fmt.Errorf("Slack API error: %s", errCheck.Error)
	}

	if respBody != nil {
		err = json.Unmarshal(respBytes, respBody)
		if err != nil {
			return errors.Wrap(err, "response body unmarshal")
		}
	}

	return nil
}
