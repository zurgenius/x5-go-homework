package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"example.com/order-service/internal/order"
)

type OrderRepository struct {
	mu    sync.RWMutex
	items []order.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{items: make([]order.Order, 0)}
}

func clone(value order.Order) order.Order {
	value.Items = append([]order.Item(nil), value.Items...)
	value.History = append([]order.StatusChange(nil), value.History...)
	return value
}

func (r *OrderRepository) Create(_ context.Context, value order.Order) (order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, err := uuid.NewV7()
	if err != nil {
		return order.Order{}, err
	}
	value.ID = id
	value.CreatedAt = time.Now().UTC()
	for i := range value.History {
		value.History[i].ChangedAt = value.CreatedAt
	}
	r.items = append(r.items, clone(value))
	return clone(value), nil
}

func (r *OrderRepository) GetByID(_ context.Context, id uuid.UUID) (order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, value := range r.items {
		if value.ID == id {
			return clone(value), nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}

func (r *OrderRepository) List(_ context.Context) ([]order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]order.Order, 0, len(r.items))
	for _, value := range r.items {
		result = append(result, clone(value))
	}
	return result, nil
}

func (r *OrderRepository) Update(_ context.Context, value order.Order) (order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, existing := range r.items {
		if existing.ID == value.ID {
			if len(value.History) > len(existing.History) {
				value.History[len(value.History)-1].ChangedAt = time.Now().UTC()
			}
			r.items[i] = clone(value)
			return clone(value), nil
		}
	}
	return order.Order{}, order.ErrOrderNotFound
}
