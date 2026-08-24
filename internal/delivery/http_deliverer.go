package delivery

import (
	"context"
	"net/http"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type HTTPDeliverer struct {
	client *http.Client
}

func NewHTTPDeliverer(client *http.Client) *HTTPDeliverer {
	return &HTTPDeliverer{
		client: client,
	}
}

func (d *HTTPDeliverer) Deliver(
	ctx context.Context,
	evt event.Event,
	sub subscription.Subscription,
) error {
	// implementation coming next
	return nil
}
