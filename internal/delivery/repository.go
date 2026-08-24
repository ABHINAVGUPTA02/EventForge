package delivery

import "context"

type Repository interface {
	Create(ctx context.Context, delivery Delivery) error
	Update(ctx context.Context, delivery Delivery) error
}
