package event

import "time"

type Event struct {
	ID string
	TenantID string
	EntityID string
	Type string
	Timestamp time.Time
	Payload []byte
}