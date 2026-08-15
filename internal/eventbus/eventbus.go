package eventbus

import (
	"context"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
)

type EventBus interface {
	Publish(ctx context.Context, event event.Event) error
	Subscribe(ctx context.Context) (<-chan event.Event, error)
}
