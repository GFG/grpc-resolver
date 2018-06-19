package resolver

import (
	"github.com/GFG/grpc-resolver/marathon"

	"google.golang.org/grpc/naming"
)

type Resolver struct {
	marathon *marathon.Client
	poller   *poll
}

// New instantiates a new resolver given a marathon uri.
func New(addr string) (*Resolver, error) {
	m := marathon.NewClient(&marathon.Config{
		URI: addr,
	})

	if err := m.Ping(); err != nil {
		return nil, err
	}

	return &Resolver{
		marathon: m,
	}, nil
}

// Resolver creates a watcher given a service name
func (r *Resolver) Resolve(name string) (naming.Watcher, error) {
	poll, err := newPoll(name, r.marathon)
	if err != nil {
		return nil, err
	}

	// Running marathon poll
	poll.run()

	return poll, nil
}

// Ready returns true if one or more probes are monitoring a grpc server
func (r *Resolver) Ready() bool {
	return len(r.poller.probes) > 0
}
