package delivery

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

func TestHTTPDelivererDeliversEvent(t *testing.T) {
	var receivedEvent event.Event

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Method != http.MethodPost {
				t.Errorf(
					"expected POST, got %s",
					r.Method,
				)
			}

			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf(
					"expected application/json, got %s",
					r.Header.Get("Content-Type"),
				)
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}

			err = json.Unmarshal(body, &receivedEvent)
			if err != nil {
				t.Fatalf("failed to decode event: %v", err)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	client := server.Client()

	deliverer := NewHTTPDeliverer(client)

	evt := event.Event{
		ID:       "evt-001",
		TenantID: "tenant-001",
		Type:     "ORDER_CREATED",
	}

	sub := subscription.Subscription{
		ID:       "sub-001",
		TenantID: "tenant-001",
		EventTypes: []string{
			"ORDER_CREATED",
		},
		Endpoint: server.URL,
	}

	err := deliverer.Deliver(
		context.Background(),
		evt,
		sub,
	)

	if err != nil {
		t.Fatalf("delivery failed: %v", err)
	}

	if receivedEvent.ID != evt.ID {
		t.Errorf(
			"expected event ID %s, got %s",
			evt.ID,
			receivedEvent.ID,
		)
	}

	if receivedEvent.Type != evt.Type {
		t.Errorf(
			"expected event type %s, got %s",
			evt.Type,
			receivedEvent.Type,
		)
	}
}

func TestHTTPDelivererReturnsErrorOnServerError(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)

	defer server.Close()

	deliverer := NewHTTPDeliverer(server.Client())

	evt := event.Event{
		ID:       "evt-002",
		TenantID: "tenant-001",
		Type:     "ORDER_CREATED",
	}

	sub := subscription.Subscription{
		ID:       "sub-002",
		TenantID: "tenant-001",
		EventTypes: []string{
			"ORDER_CREATED",
		},
		Endpoint: server.URL,
	}

	err := deliverer.Deliver(
		context.Background(),
		evt,
		sub,
	)

	if err == nil {
		t.Fatal("expected delivery to fail")
	}
}
