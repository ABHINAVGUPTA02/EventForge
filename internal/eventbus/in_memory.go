package eventbus

import (
	"context"
	"sync"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
)

const subscriberBufferSize = 100

type InMemoryBus struct {
	mu          sync.RWMutex
	subscribers []chan event.Event
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		subscribers: make([]chan event.Event, 0),
	}
}

func (b *InMemoryBus) Subscribe(ctx context.Context) (<-chan event.Event, error) {
	ch := make(chan event.Event, subscriberBufferSize)

	b.mu.Lock()
	b.subscribers = append(b.subscribers, ch)
	b.mu.Unlock()

	return ch, nil
}

func (b *InMemoryBus) Publish(ctx context.Context, evnt event.Event) error {
	b.mu.RLock()
	subscribers := make([]chan event.Event, len(b.subscribers))
	copy(subscribers, b.subscribers)
	b.mu.RUnlock()

	for _, ch := range subscribers {
		select {
		case ch <- evnt:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
