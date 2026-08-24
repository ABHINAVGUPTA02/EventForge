package delivery

import (
	"context"
	"sync"

	"github.com/ABHINAVGUPTA02/EventForge/internal/event"
	"github.com/ABHINAVGUPTA02/EventForge/internal/subscription"
)

type Manager struct {
	deliverer   Deliverer
	tasks       chan DeliveryTask
	workerCount int
	repository  Repository

	wg sync.WaitGroup
}

func NewManager(deliverer Deliverer, repository Repository, workerCount int, queueSize int) *Manager {
	return &Manager{
		deliverer:   deliverer,
		repository:  repository,
		tasks:       make(chan DeliveryTask, queueSize),
		workerCount: workerCount,
	}
}

func (m *Manager) Start(ctx context.Context) {
	for i := 0; i < m.workerCount; i++ {
		m.wg.Add(1)

		go m.worker(ctx)
	}
}

func (m *Manager) worker(ctx context.Context) {
	defer m.wg.Done()

	for {
		select {
		case task, ok := <-m.tasks:
			if !ok {
				return
			}

			task.Delivery.Status = StatusDelivering
			task.Delivery.Attempt++

			if err := m.repository.Update(
				ctx,
				task.Delivery,
			); err != nil {
				continue
			}

			err := m.deliverer.Deliver(
				ctx,
				task.Event,
				task.Subscription,
			)

			if err != nil {
				task.Delivery.Status = StatusRetrying
			} else {
				task.Delivery.Status = StatusDelivered
			}

			if err := m.repository.Update(
				ctx,
				task.Delivery,
			); err != nil {
				continue
			}

		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Submit(
	ctx context.Context,
	task DeliveryTask,
) error {
	select {
	case m.tasks <- task:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Manager) Stop() {
	close(m.tasks)
	m.wg.Wait()
}

func (m *Manager) Enqueue(
	ctx context.Context,
	evt event.Event,
	sub subscription.Subscription,
) error {
	delivery := Delivery{
		ID:             evt.ID + ":" + sub.ID,
		EventID:        evt.ID,
		SubscriptionID: sub.ID,
		Status:         StatusPending,
		Attempt:        0,
	}

	if err := m.repository.Create(ctx, delivery); err != nil {
		return err
	}

	task := DeliveryTask{
		Event:        evt,
		Subscription: sub,
		Delivery:     delivery,
	}

	return m.Submit(ctx, task)
}
