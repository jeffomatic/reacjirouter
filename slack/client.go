package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const defaultURLPrefix = "https://slack.com/api"

type Client struct {
	URLPrefix   string
	AccessToken string
}

type APIError string

func (e APIError) Error() string {
	return string(e)
}

func GetAPIError(err error) *APIError {
	if apiErr, ok := err.(APIError); ok {
		return &apiErr
	}

	return nil
}

func NewClient(token string) *Client {
	return &Client{URLPrefix: defaultURLPrefix, AccessToken: token}
}

func (c *Client) Call(method Method, body interface{}, returnData interface{}) error {
	methodInfo := infoByMethod[method]

	reqBody := []byte{}
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return errors.Wrap(err, "json encode request body")
		}

		if methodInfo.formEncodedRequest {
			// We use JSON struct tagging for serialization. It's not great but better
			// than inventing a new tagging system just for form-encoded variables.
			var params map[string]string
			err = json.Unmarshal(reqBody, &params)
			if err != nil {
				return errors.Wrap(err, "invalid request body")
			}

			formVars := make(url.Values)
			for k, v := range params {
				formVars[k] = []string{v}
			}

			reqBody = []byte(formVars.Encode())
		}
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.URLPrefix+"/"+string(method),
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return errors.Wrap(err, "HTTP request creation")
	}

	if methodInfo.formEncodedRequest {
		req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	} else {
		req.Header["Content-Type"] = []string{"application/json"}
	}

	if !methodInfo.noAccessToken {
		req.Header["Authorization"] = []string{"Bearer " + c.AccessToken}
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return errors.Wrap(err, "executing HTTP request")
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "reading HTTP response body")
	}

	errCheck := struct {
		Ok    bool
		Error APIError
	}{}
	err = json.Unmarshal(respBytes, &errCheck)
	if err != nil {
		return errors.Wrap(err, "response body error-check decode")
	}
	if !errCheck.Ok {
		return errCheck.Error
	}

	if returnData != nil {
		err = json.Unmarshal(respBytes, returnData)
		if err != nil {
			return errors.Wrap(err, "response body unmarshal")
		}
	}

	return nil
}
