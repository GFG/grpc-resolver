package marathon

import "strconv"

type ApplicationList struct {
	Apps []*Application `json:"apps"`
}

// Application represents the object for an application in marathon
type Application struct {
	ID        string             `json:"id,omitempty"`
	Container *Container         `json:"container,omitempty"`
	Labels    *map[string]string `json:"labels,omitempty"`
}

// Container is the definition for a container type in marathon
type Container struct {
	Type   string  `json:"type,omitempty"`
	Docker *Docker `json:"docker,omitempty"`
}

// Docker is the docker definition from a marathon application
type Docker struct {
	PortMappings *[]PortMapping `json:"portMappings,omitempty"`
}

// PortMapping is the portmapping structure between container and mesos
type PortMapping struct {
	ContainerPort int                `json:"containerPort,omitempty"`
	HostPort      int                `json:"hostPort"`
	Labels        *map[string]string `json:"labels,omitempty"`
	Name          string             `json:"name,omitempty"`
	ServicePort   int                `json:"servicePort,omitempty"`
	Protocol      string             `json:"protocol,omitempty"`
}

// PortDefinition represents a port that should be considered part of
// a resource. Port definitions are necessary when you are using HOST
// networking
type PortDefinition struct {
	Port     *int               `json:"port,omitempty"`
	Protocol string             `json:"protocol,omitempty"`
	Name     string             `json:"name,omitempty"`
	Labels   *map[string]string `json:"labels,omitempty"`
}

type TaskList struct {
	Tasks []*Task `json:"tasks,omitempty"`
}
// Task represents the definition for a marathon task
type Task struct {
	ID    string `json:"id"`
	AppID string `json:"appId"`
	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

// Addr returns the task full address given a port index (format: 192.168.0.1:8080)
func (t *Task) Addr(portIndex int64) string {
	return t.Host + ":" + strconv.FormatInt(int64(t.Ports[portIndex]), 10)
}

// IPAddress represents a task's IP address and protocol.
type IPAddress struct {
	IPAddress string `json:"ipAddress"`
	Protocol  string `json:"protocol"`
}
