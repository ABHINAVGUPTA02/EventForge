package router

import (
	"sync"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
)

type Router struct {
	mu            sync.RWMutex
	subscriptions map[string]Subscription
}

func NewRouter() *Router {
	return &Router{
		subscriptions: make(map[string]Subscription),
	}
}

func (r *Router) Subscribe(sub Subscription) {
	r.mu.Lock()
	defer r.mu.RUnlock()

	r.subscriptions[sub.ID] = sub
}

func (r *Router) Match(evnt event.Event) []Subscription {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []Subscription

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
