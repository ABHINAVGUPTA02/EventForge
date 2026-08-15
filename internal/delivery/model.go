package delivery

type Status string

const (
	StatusPending      Status = "PENDING"
	StatusDelivering   Status = "DELIVERING"
	StatusDelivered    Status = "DELIVERED"
	StatusRetrying     Status = "RETRYING"
	StatusDeadLettered Status = "DEAD_LETTERED"
)

type Delivery struct {
	ID             string
	EventID        string
	SubscriptionID string
	Status         string
	Attempt        int
}
