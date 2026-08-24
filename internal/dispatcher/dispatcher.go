package dispatcher

import (
	"context"

	"github.com/ABHINAVGUPTA02/EventForge/internal/delivery"
	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/router"
)

type Dispatcher struct {
	router  *router.Router
	manager *delivery.Manager
}

func NewDispatcher(
	router *router.Router,
	manager *delivery.Manager,
) *Dispatcher {
	return &Dispatcher{
		router:  router,
		manager: manager,
	}
}

func (d *Dispatcher) Dispatch(
	ctx context.Context,
	evt event.Event,
) error {
	subscriptions := d.router.Match(evt)

	for _, sub := range subscriptions {
		if err := d.manager.Enqueue(ctx, evt, sub); err != nil {
			return err
		}
	}

	return nil
}
