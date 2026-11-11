package service_test

import (
	"context"
	"errors"
	"math"
	"testing"

	"example.com/order-service/internal/order"
	"example.com/order-service/internal/repository/memory"
	"example.com/order-service/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type fakeRepository struct {
	order.Repository
	create func(context.Context, order.Order) (order.Order, error)
}

func (f fakeRepository) Create(ctx context.Context, value order.Order) (order.Order, error) {
	return f.create(ctx, value)
}

func TestCreate(t *testing.T) {
	valid := order.Item{Name: "Keyboard", Quantity: 2, Price: 100}
	cases := []struct {
		name  string
		items []order.Item
		valid bool
	}{
		{"valid", []order.Item{valid, {Name: "Mouse", Quantity: 1, Price: 50}}, true},
		{"empty items", nil, false},
		{"empty name", []order.Item{{Name: " ", Quantity: 1, Price: 1}}, false},
		{"zero quantity", []order.Item{{Name: "x", Quantity: 0, Price: 1}}, false},
		{"negative quantity", []order.Item{{Name: "x", Quantity: -1, Price: 1}}, false},
		{"zero price", []order.Item{{Name: "x", Quantity: 1, Price: 0}}, false},
		{"negative price", []order.Item{{Name: "x", Quantity: 1, Price: -1}}, false},
		{"overflow", []order.Item{{Name: "x", Quantity: 2, Price: math.MaxInt64}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			repo := fakeRepository{create: func(_ context.Context, value order.Order) (order.Order, error) {
				called = true
				return value, nil
			}}
			got, err := service.NewOrderService(repo).Create(context.Background(), order.CreateOrderInput{Items: tc.items})
			if !tc.valid {
				require.ErrorIs(t, err, order.ErrInvalidOrder)
				require.False(t, called)
				return
			}
			require.NoError(t, err)
			require.True(t, called)
			require.EqualValues(t, 250, got.TotalPrice)
			require.Equal(t, order.StatusNew, got.Status)
			require.Len(t, got.History, 1)
			require.True(t, got.CreatedAt.IsZero())
			require.True(t, got.History[0].ChangedAt.IsZero())
		})
	}
}

func TestOrderLifecycle(t *testing.T) {
	ctx := context.Background()
	svc := service.NewOrderService(memory.NewOrderRepository())
	created, err := svc.Create(ctx, order.CreateOrderInput{Items: []order.Item{{Name: "item", Quantity: 1, Price: 2}}})
	require.NoError(t, err)
	require.Equal(t, uuid.Version(7), created.ID.Version())
	require.False(t, created.CreatedAt.IsZero())
	require.Equal(t, created.CreatedAt, created.History[0].ChangedAt)
	got, err := svc.GetByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, created, got)
	_, err = svc.GetByID(ctx, uuid.New())
	require.ErrorIs(t, err, order.ErrOrderNotFound)
	list, err := svc.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 1)
	_, err = svc.UpdateStatus(ctx, created.ID, order.StatusProcessing)
	require.NoError(t, err)
	_, err = svc.UpdateStatus(ctx, created.ID, order.StatusCompleted)
	require.NoError(t, err)
	history, err := svc.History(ctx, created.ID)
	require.NoError(t, err)
	require.Len(t, history, 3)
	for _, change := range history {
		require.False(t, change.ChangedAt.IsZero())
	}
	require.Equal(t, []order.Status{order.StatusNew, order.StatusProcessing, order.StatusCompleted}, []order.Status{history[0].Status, history[1].Status, history[2].Status})
	for i := 1; i < len(history); i++ {
		require.False(t, history[i].ChangedAt.Before(history[i-1].ChangedAt))
	}
	_, err = svc.History(ctx, uuid.New())
	require.ErrorIs(t, err, order.ErrOrderNotFound)
}

func TestStatusTransitions(t *testing.T) {
	cases := []struct {
		name  string
		steps []order.Status
		want  error
	}{
		{"new to processing", []order.Status{order.StatusProcessing}, nil},
		{"new to cancelled", []order.Status{order.StatusCancelled}, nil},
		{"processing to completed", []order.Status{order.StatusProcessing, order.StatusCompleted}, nil},
		{"processing to cancelled", []order.Status{order.StatusProcessing, order.StatusCancelled}, nil},
		{"new to completed", []order.Status{order.StatusCompleted}, order.ErrInvalidStatusTransition},
		{"same status", []order.Status{order.StatusNew}, order.ErrInvalidStatusTransition},
		{"completed is final", []order.Status{order.StatusProcessing, order.StatusCompleted, order.StatusCancelled}, order.ErrInvalidStatusTransition},
		{"cancelled is final", []order.Status{order.StatusCancelled, order.StatusProcessing}, order.ErrInvalidStatusTransition},
		{"unknown", []order.Status{"bad"}, order.ErrInvalidStatus},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			svc := service.NewOrderService(memory.NewOrderRepository())
			value, err := svc.Create(ctx, order.CreateOrderInput{Items: []order.Item{{Name: "x", Quantity: 1, Price: 1}}})
			require.NoError(t, err)
			for i, status := range tc.steps {
				_, err = svc.UpdateStatus(ctx, value.ID, status)
				if i == len(tc.steps)-1 && tc.want != nil {
					require.ErrorIs(t, err, tc.want)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestRepositoryErrors(t *testing.T) {
	root := errors.New("storage unavailable")
	repo := fakeRepository{create: func(context.Context, order.Order) (order.Order, error) { return order.Order{}, root }}
	_, err := service.NewOrderService(repo).Create(context.Background(), order.CreateOrderInput{Items: []order.Item{{Name: "x", Quantity: 1, Price: 1}}})
	require.ErrorIs(t, err, root)
	_, err = service.NewOrderService(memory.NewOrderRepository()).UpdateStatus(context.Background(), uuid.New(), order.StatusProcessing)
	require.ErrorIs(t, err, order.ErrOrderNotFound)
}
