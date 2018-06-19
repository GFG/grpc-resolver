package resolver

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Probe struct {
	addr   string
	conn   *grpc.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

func newProbe(addr string, timeout time.Duration) (*Probe, error) {
	ctx, cancel := context.WithCancel(context.Background())

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Probe{
		addr:   addr,
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (p *Probe) exec() chan connectivity.State {
	out := make(chan connectivity.State)

	go func() {
		defer close(out)
		for {
			current := p.conn.GetState()

			out <- current

			ok := p.conn.WaitForStateChange(p.ctx, current)
			if !ok {
				if p.ctx.Err() == context.DeadlineExceeded {
					out <- connectivity.TransientFailure
				}
				return
			}
		}
	}()

	return out
}

func (p *Probe) close() {
	p.cancel()
	if p.conn != nil {
		_ = p.conn.Close()
	}
}
