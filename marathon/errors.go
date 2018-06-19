package marathon

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Error is a trivial implementation of error
type Error struct {
	message string `json:"message"`
}

// Error returns the description of the error
func (e *Error) Error() string {
	return e.message
}

func parseError(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = &Error{}
	if err := json.Unmarshal(b, err); err != nil {
		return err
	}

	return err
}
