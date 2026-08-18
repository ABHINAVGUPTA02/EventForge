package router

import (
	"sync"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type Router struct {
	mu            sync.RWMutex
	subscriptions map[string]subscription.Subscription
}

func NewRouter() *Router {
	return &Router{
		subscriptions: make(map[string]subscription.Subscription),
	}
}

func (r *Router) Subscribe(sub subscription.Subscription) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscriptions[sub.ID] = sub
}

func (r *Router) Match(evnt event.Event) []subscription.Subscription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []subscription.Subscription

	for _, sub := range r.subscriptions {
		if sub.TenantID != evnt.TenantID {
			continue
		}

		for _, eventType := range sub.EventTypes {
			if eventType == evnt.Type {
				matches = append(matches, sub)
				break
			}
		}
	}

	return matches
}
