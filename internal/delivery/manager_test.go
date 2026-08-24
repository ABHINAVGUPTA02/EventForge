package delivery

import (
	"context"
	"sync"
	"testing"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
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

func TestManagerProcessesTasks(t *testing.T) {
	ctx := context.Background()

	deliverer := &fakeDeliverer{}
	repo := NewInMemoryRepository()

	manager := NewManager(
		deliverer,
		repo,
		3,
		10,
	)

	manager.Start(ctx)

	task := DeliveryTask{
		Event: event.Event{
			ID:       "evt-001",
			TenantID: "tenant-001",
			Type:     "ORDER_CREATED",
		},
		Subscription: subscription.Subscription{
			ID:       "sub-001",
			TenantID: "tenant-001",
			EventTypes: []string{
				"ORDER_CREATED",
			},
			Endpoint: "http://example.com",
		},
		Delivery: Delivery{
			ID:             "delivery-001",
			EventID:        "evt-001",
			SubscriptionID: "sub-001",
			Status:         StatusPending,
			Attempt:        0,
		},
	}

	for i := 0; i < 5; i++ {
		err := manager.Submit(ctx, task)
		if err != nil {
			t.Fatalf("failed to submit task: %v", err)
		}
	}

	manager.Stop()

	if got := deliverer.Calls(); got != 5 {
		t.Fatalf("expected 5 deliveries, got %d", got)
	}
}

func TestManagerUpdatesDeliveryStatus(t *testing.T) {
	ctx := context.Background()

	deliverer := &fakeDeliverer{}
	repo := NewInMemoryRepository()

	manager := NewManager(
		deliverer,
		repo,
		1,
		10,
	)

	manager.Start(ctx)

	task := DeliveryTask{
		Event: event.Event{
			ID:       "evt-002",
			TenantID: "tenant-001",
			Type:     "ORDER_CREATED",
		},
		Subscription: subscription.Subscription{
			ID:       "sub-002",
			TenantID: "tenant-001",
			EventTypes: []string{
				"ORDER_CREATED",
			},
			Endpoint: "http://example.com",
		},
		Delivery: Delivery{
			ID:             "delivery-002",
			EventID:        "evt-002",
			SubscriptionID: "sub-002",
			Status:         StatusPending,
			Attempt:        0,
		},
	}

	err := manager.Submit(ctx, task)
	if err != nil {
		t.Fatalf("failed to submit task: %v", err)
	}

	manager.Stop()

	got, ok := repo.Get("delivery-002")
	if !ok {
		t.Fatal("delivery was not found in repository")
	}

	if got.Status != StatusDelivered {
		t.Fatalf(
			"expected status %s, got %s",
			StatusDelivered,
			got.Status,
		)
	}

	if got.Attempt != 1 {
		t.Fatalf(
			"expected attempt 1, got %d",
			got.Attempt,
		)
	}
}
