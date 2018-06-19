package marathon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) apiCall(method, path string, reader io.Reader, result interface{}) error {
	req, err := c.makeRequest(method, path, reader)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return parseError(resp)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); nil != err {
			println(err.Error())
			return err
		}
	}

	return nil
}

func (c *Client) makeRequest(method, path string, reader io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.config.URI, path)

	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	if c.config.HTTPBasicAuthUser != "" && c.config.HTTPBasicAuthPassword != "" {
		request.SetBasicAuth(c.config.HTTPBasicAuthUser, c.config.HTTPBasicAuthPassword)
	}

	if c.config.DCOSToken != "" {
		request.Header.Add("Authorization", "token="+c.config.HTTPBasicAuthUser)
	}

	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
