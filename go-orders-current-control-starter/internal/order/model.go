package order

import (
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

type Item struct {
	Name     string
	Quantity int
	Price    int64
}

type StatusChange struct {
	Status    Status
	ChangedAt time.Time
}

type Order struct {
	ID         uuid.UUID
	Items      []Item
	TotalPrice int64
	Status     Status
	CreatedAt  time.Time
	History    []StatusChange
}

type CreateOrderInput struct {
	Items []Item
}

type UpdateStatusInput struct {
	Status Status
}
