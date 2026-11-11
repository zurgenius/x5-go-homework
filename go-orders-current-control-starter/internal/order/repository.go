package order

import (
	"context"
	"github.com/google/uuid"
)

// Repository describes the storage operations required by the business layer.
// Create assigns the ID and initial timestamps. Update assigns the timestamp
// to a newly appended status change.
type Repository interface {
	Create(ctx context.Context, value Order) (Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (Order, error)
	List(ctx context.Context) ([]Order, error)
	Update(ctx context.Context, value Order) (Order, error)
}
