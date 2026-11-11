package service

import (
	"context"
	"math"
	"strings"
	"sync"

	"github.com/google/uuid"

	"example.com/order-service/internal/order"
)

type OrderService struct {
	repository order.Repository
	mu         sync.Mutex
}

func NewOrderService(repository order.Repository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Create(ctx context.Context, input order.CreateOrderInput) (order.Order, error) {
	if len(input.Items) == 0 {
		return order.Order{}, order.ErrInvalidOrder
	}
	var total int64
	for _, item := range input.Items {
		if strings.TrimSpace(item.Name) == "" || item.Quantity <= 0 || item.Price <= 0 || item.Price > (math.MaxInt64-total)/int64(item.Quantity) {
			return order.Order{}, order.ErrInvalidOrder
		}
		total += item.Price * int64(item.Quantity)
	}
	return s.repository.Create(ctx, order.Order{
		Items: input.Items, TotalPrice: total, Status: order.StatusNew,
		History: []order.StatusChange{{Status: order.StatusNew}},
	})
}

func (s *OrderService) GetByID(ctx context.Context, id uuid.UUID) (order.Order, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *OrderService) List(ctx context.Context) ([]order.Order, error) {
	return s.repository.List(ctx)
}

func (s *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, status order.Status) (order.Order, error) {
	switch status {
	case order.StatusNew, order.StatusProcessing, order.StatusCompleted, order.StatusCancelled:
	default:
		return order.Order{}, order.ErrInvalidStatus
	}
	// ponytail: one lock serializes updates in this process; move transitions into
	// transactional storage if multiple service instances need to update orders.
	s.mu.Lock()
	defer s.mu.Unlock()
	value, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return order.Order{}, err
	}
	if !canTransition(value.Status, status) {
		return order.Order{}, order.ErrInvalidStatusTransition
	}
	value.Status = status
	value.History = append(value.History, order.StatusChange{Status: status})
	return s.repository.Update(ctx, value)
}

func (s *OrderService) History(ctx context.Context, id uuid.UUID) ([]order.StatusChange, error) {
	value, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return value.History, nil
}

func canTransition(from, to order.Status) bool {
	switch from {
	case order.StatusNew:
		return to == order.StatusProcessing || to == order.StatusCancelled
	case order.StatusProcessing:
		return to == order.StatusCompleted || to == order.StatusCancelled
	default:
		return false
	}
}
