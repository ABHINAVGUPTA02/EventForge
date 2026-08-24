package dispatcher_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ABHINAVGUPTA02/EventForge/internal/delivery"
	"github.com/ABHINAVGUPTA02/EventForge/internal/dispatcher"
	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/router"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

func TestDispatcherDeliversEventEndToEnd(t *testing.T) {
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

			if err := json.NewDecoder(r.Body).Decode(&receivedEvent); err != nil {
				t.Fatalf(
					"failed to decode event: %v",
					err,
				)
			}

			w.WriteHeader(http.StatusOK)
		}),
	)

	defer server.Close()

	repo := delivery.NewInMemoryRepository()

	httpClient := server.Client()

	deliverer := delivery.NewHTTPDeliverer(httpClient)

	manager := delivery.NewManager(
		deliverer,
		repo,
		3,
		10,
	)

	r := router.NewRouter()

	d := dispatcher.NewDispatcher(
		r,
		manager,
	)

	r.Subscribe(subscription.Subscription{
		ID:       "sub-001",
		TenantID: "tenant-001",
		EventTypes: []string{
			"ORDER_CREATED",
		},
		Endpoint: server.URL,
	})

	ctx := context.Background()

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
		t.Fatalf(
			"dispatch failed: %v",
			err,
		)
	}

	manager.Stop()

	if receivedEvent.ID != "evt-001" {
		t.Errorf(
			"expected event ID evt-001, got %s",
			receivedEvent.ID,
		)
	}

	if receivedEvent.Type != "ORDER_CREATED" {
		t.Errorf(
			"expected event type ORDER_CREATED, got %s",
			receivedEvent.Type,
		)
	}
}
