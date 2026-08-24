package delivery

import (
	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type DeliveryTask struct {
	Event        event.Event
	Subscription subscription.Subscription
	Delivery     Delivery
}
