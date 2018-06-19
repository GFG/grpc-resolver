package resolver

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GFG/grpc-resolver/marathon"

	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/naming"
)

type poll struct {
	label      string
	appID      string
	portIndex  int64
	probes     map[string]*Probe
	marathon   *marathon.Client
	updates    chan []*marathon.Task
	unregister chan string
	done       chan bool
}

func newPoll(label string, m *marathon.Client) (*poll, error) {
	apps, err := m.Applications(label)
	if err != nil {
		return nil, err
	}

	if len(apps) != 1 {
		return nil, errors.New("Duplicate labels or label not found")
	}

	// Parsing service port index
	portIndex := int64(-1)
	for k := range *apps[0].Labels {
		if k == label {
			res := strings.Split(k, "_")
			if len(res) != 3 {
				continue
			}

			portIndex, err = strconv.ParseInt(res[1], 10, 64)
			if err != nil {
				return nil, errors.New("Failed parse port index in tag")
			}
		}
	}

	if portIndex < 0 {
		return nil, errors.New("label not found")
	}

	return &poll{
		label:      label,
		portIndex:  portIndex,
		appID:      apps[0].ID,
		probes:     make(map[string]*Probe),
		updates:    make(chan []*marathon.Task, 0),
		unregister: make(chan string, 0),
		done:       make(chan bool, 1),
		marathon:   m,
	}, nil
}

// poll polls the service tasks' states in marathon
func (p *poll) poll() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ticker.C:
			tasks, err := p.marathon.Tasks(p.appID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "couldn't retrieve tasks in marathon: %v. Trying again...", err)
				continue
			}

			p.updates <- tasks
		case <-p.done:
			return
		}
	}
}

// Next blocks until an update or error happens in the service polling
func (p *poll) Next() ([]*naming.Update, error) {
	var ups []*naming.Update

	for {
		select {
		case tasks, ok := <-p.updates:
			if !ok {
				p.updates = nil
				return nil, errors.New("poller unrecoverable")
			}

			for _, task := range tasks {
				if _, ok := p.probes[task.Addr(p.portIndex)]; ok {
					// If the task is already registered,
					// there is nothing to do.
					continue
				}

				go p.processTask(task)

				ups = append(ups, &naming.Update{
					Addr: task.Addr(p.portIndex),
					Op:   naming.Add,
				})
			}

			if len(ups) > 0 {
				return ups, nil
			}
		case addr, ok := <-p.unregister:
			if !ok {
				p.unregister = nil
				return nil, errors.New("poller unrecoverable")
			}

			delete(p.probes, addr)

			ups = append(ups, &naming.Update{
				Addr: addr,
				Op:   naming.Delete,
			})

			return ups, nil
		case <-p.done:
			return ups, nil
		}
	}
}

func (p *poll) processTask(task *marathon.Task) {
	probe, err := newProbe(task.Addr(p.portIndex), time.Second*5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to instantiate probe: %v", err)
		p.unregister <- probe.addr
		return
	}

	p.probes[probe.addr] = probe

	out := probe.exec()

	for {
		select {
		case <-p.done:
			probe.close()
			return
		case state := <-out:
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				p.unregister <- probe.addr
				probe.close()
				return
			}
		}
	}

}

func (p *poll) run() {
	go p.poll()
}

// Close closes the polling and the probes monitoring
func (p *poll) Close() {
	close(p.done)
	close(p.unregister)
	close(p.updates)
}
