package delivery

import (
	"context"
	"sync"
)

type InMemoryRepository struct {
	mu         sync.RWMutex
	deliveries map[string]Delivery
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		deliveries: make(map[string]Delivery),
	}
}

func (r *InMemoryRepository) Create(
	ctx context.Context,
	delivery Delivery,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.deliveries[delivery.ID] = delivery

	return nil
}

func (r *InMemoryRepository) Update(
	ctx context.Context,
	delivery Delivery,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.deliveries[delivery.ID] = delivery

	return nil
}

func (r *InMemoryRepository) Get(id string) (Delivery, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	delivery, ok := r.deliveries[id]
	return delivery, ok
}
