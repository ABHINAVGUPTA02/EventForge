package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		sub.Endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"delivery failed with status code %d",
			resp.StatusCode,
		)
	}

	return nil
}
