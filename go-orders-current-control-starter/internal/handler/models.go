package handler

import (
	"time"

	"example.com/order-service/internal/order"
	"github.com/google/uuid"
)

type itemRequest struct {
	Name     string `json:"name" validate:"required,min=1"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Price    int64  `json:"price" validate:"required,min=1"`
}

type createOrderRequest struct {
	Items []itemRequest `json:"items" validate:"required,min=1"`
}

type updateStatusRequest struct {
	Status order.Status `json:"status" validate:"required,oneof=new processing completed cancelled"`
}

type itemResponse struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price"`
}

type orderResponse struct {
	ID         uuid.UUID      `json:"id" format:"uuid"`
	Items      []itemResponse `json:"items"`
	TotalPrice int64          `json:"total_price"`
	Status     order.Status   `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
}

type statusChangeResponse struct {
	Status    order.Status `json:"status"`
	ChangedAt time.Time    `json:"changed_at"`
}

type errorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func toOrderResponse(value order.Order) orderResponse {
	items := make([]itemResponse, 0, len(value.Items))
	for _, item := range value.Items {
		items = append(items, itemResponse{Name: item.Name, Quantity: item.Quantity, Price: item.Price})
	}
	return orderResponse{ID: value.ID, Items: items, TotalPrice: value.TotalPrice, Status: value.Status, CreatedAt: value.CreatedAt}
}
