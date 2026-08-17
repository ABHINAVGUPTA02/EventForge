package eventbus

import (
	"context"
	"testing"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"

	"time"
)

//Test Publish and Subscribe
func TestPublishAndSubscibe(t *testing.T) {
	//context
	ctx := context.Background()
	//creating a new in memory bus
	bus := NewInMemoryBus()

	//subscribe to a channel
	ch, err := bus.Subscribe(ctx)
	if err != nil {
		t.Fatalf("failed to subscribe: %v", err)
	}

	//creating a dummy event
	dummyEvent := event.Event{
		ID:        "evt-001",
		TenantID:  "tenant-001",
		Type:      "ORDER_CREATED",
		EntityID:  "order-101",
		Timestamp: time.Now(),
		Payload:   []byte(`{"order_id":"order-101"}`),
	}

	//publish to a channel
	err = bus.Publish(ctx, dummyEvent)
	if err != nil {
		t.Fatalf("failed to publish event: %v", err)
	}

	//receiving a channel
	received := <-ch

	//assert condition
	if received.ID != dummyEvent.ID {
		t.Errorf("expected event ID %s, got %s", dummyEvent.ID, received.ID)
	}
}