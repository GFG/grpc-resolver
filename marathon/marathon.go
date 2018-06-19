package marathon

import "fmt"

type Client struct {
	config *Config
}

// Config represents the marathon client configuration object
type Config struct {
	HTTPBasicAuthUser     string
	HTTPBasicAuthPassword string
	DCOSToken             string
	URI                   string
}

// NewClient instantiates a new marathon client
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

// Applications returns a set of applications according to a label
func (c *Client) Applications(label string) ([]*Application, error) {
	apps := &ApplicationList{
		Apps: make([]*Application, 0),
	}

	path := fmt.Sprintf("/v2/apps?label=%s", label)

	if err := c.apiCall("GET", path, nil, &apps); err != nil {
		return nil, err
	}

	return apps.Apps, nil
}

// Tasks returns a specific application's set of tasks
func (c *Client) Tasks(appID string) ([]*Task, error) {
	tasks := &TaskList{
		Tasks: make([]*Task, 0),
	}

	if err := c.apiCall(
		"GET",
		fmt.Sprintf("/v2/apps/%s/tasks", appID),
		nil,
		&tasks); err != nil {
		return nil, err
	}

	return tasks.Tasks, nil
}

// Ping returns an error if the marathon framework is unreachable
func (c *Client) Ping() error {
	return c.apiCall("GET", "/ping", nil, nil)
}

// URI returns the marathon uri
func (c *Client) URI() string {
	return c.config.URI
}
