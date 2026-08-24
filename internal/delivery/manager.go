package delivery

import (
	"context"
	"sync"
)

type Manager struct {
	deliverer   Deliverer
	tasks       chan DeliveryTask
	workerCount int

	wg sync.WaitGroup
}

func NewManager(deliverer Deliverer, workerCount int, queueSize int) *Manager {
	return &Manager{
		deliverer:   deliverer,
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

			_ = m.deliverer.Deliver(
				ctx,
				task.Event,
				task.Subscription,
			)
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
