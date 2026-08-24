package delivery

import (
	"context"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type Deliverer interface {
	Deliver(
		ctx context.Context,
		evt event.Event,
		sub subscription.Subscription,
	) error
}
