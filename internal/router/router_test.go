package router

import (
	"testing"
	"time"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

func TestRouterSubscribeAndMatching(t *testing.T) {
	r := NewRouter()

	sub := subscription.Subscription{
		ID:             "001",
		TenantID:       "T-001",
		EventTypes:     []string{"ORDER_CREATED"},
		Endpoint:       "https://test-endpoint.com",
		MaxConcurrency: 1,
		RateLimit:      1,
	}

	r.Subscribe(sub)

	evnt := event.Event{
		ID:        "001-001",
		TenantID:  "T-001",
		EntityID:  "E-001",
		Type:      "ORDER_CREATED",
		Timestamp: time.Now(),
		Payload:   make([]byte, 0),
	}

	match := r.Match(evnt)

	if len(match) != 1 {
		t.Fatalf("Error in receiving a event for a subscription %s", sub.ID)
	}

	if match[0].TenantID != evnt.TenantID {
		t.Fatalf("Error in the received subscription as it does not match the TenantID %s", match[0].TenantID)
	}
}
