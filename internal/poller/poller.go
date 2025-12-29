package poller

import (
	"context"
	"fmt"
	"go-worker/internal/example"
	"go-worker/internal/poller/dispatcher"
	"time"

	"go.uber.org/fx"
)

/*
Fx App
 ├── Dispatcher (long-lived)
 │    └── Worker pools
 ├── Poller (long-lived background process)
 │    └── Infinite loop
 │         ├── Check condition (DB / cache / API)
 │         └── Dispatch job(s)
 └── Servers
*/

type Poller struct {
	dispatcher *dispatcher.Service
	interval   time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(
	dispatcher *dispatcher.Service,
) *Poller {
	return &Poller{
		dispatcher: dispatcher,
		interval:   2 * time.Second, // configurable
	}
}

func (p *Poller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		fmt.Println("poller For loop ...")
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			p.tick(ctx)
		}
	}
}

func (p *Poller) tick(ctx context.Context) {
	// Example condition (replace with DB/cache logic)
	fmt.Println("tick running ...")
	if !p.conditionMet(ctx) {
		return
	}

	j := example.NewExampleJob("job-123", "email")

	_ = p.dispatcher.Dispatch(j)
}

func (p *Poller) conditionMet(ctx context.Context) bool {
	// DB / Redis / API / feature flag / backlog size
	return true
}

func RegisterLifecycle(lc fx.Lifecycle, p *Poller) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("starting dispatcher...")
			p.ctx, p.cancel = context.WithCancel(context.Background())
			go p.Run(p.ctx)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("stopping poller...")
			if p.cancel != nil {
				p.cancel()
			}
			return nil
		},
	})
}
