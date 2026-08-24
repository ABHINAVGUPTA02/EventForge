package dispatcher

import (
	"context"
	"sync"
	"testing"

	"github.com/ABHINAVGUPTA02/EventForge/internal/delivery"
	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/router"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type fakeDeliverer struct {
	mu    sync.Mutex
	calls int
}

func (d *fakeDeliverer) Deliver(
	ctx context.Context,
	evt event.Event,
	sub subscription.Subscription,
) error {
	d.mu.Lock()
	d.calls++
	d.mu.Unlock()

	return nil
}

func (d *fakeDeliverer) Calls() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.calls
}

func TestDispatcherDeliversEvent(t *testing.T) {
	ctx := context.Background()

	deliverer := &fakeDeliverer{}
	repo := delivery.NewInMemoryRepository()

	manager := delivery.NewManager(
		deliverer,
		repo,
		3,
		10,
	)

	r := router.NewRouter()

	d := NewDispatcher(
		r,
		manager,
	)

	r.Subscribe(subscription.Subscription{
		ID:       "sub-001",
		TenantID: "tenant-001",
		EventTypes: []string{
			"ORDER_CREATED",
		},
		Endpoint: "http://example.com",
	})

	manager.Start(ctx)

	err := d.Dispatch(
		ctx,
		event.Event{
			ID:       "evt-001",
			TenantID: "tenant-001",
			Type:     "ORDER_CREATED",
		},
	)

	if err != nil {
		t.Fatalf("dispatch failed: %v", err)
	}

	manager.Stop()

	if got := deliverer.Calls(); got != 1 {
		t.Fatalf("expected 1 delivery, got %d", got)
	}
}
