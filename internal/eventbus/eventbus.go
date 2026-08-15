package eventbus

import (
	"context"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
)

type EventBus interface {
	publish(ctx context.Context, event event.Event) error
	subscribe(ctx context.Context) (<-chan event.Event, error)
}
