package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	_ "example.com/order-service/docs"
	"example.com/order-service/internal/order"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type OrderService interface {
	Create(context.Context, order.CreateOrderInput) (order.Order, error)
	GetByID(context.Context, uuid.UUID) (order.Order, error)
	List(context.Context) ([]order.Order, error)
	UpdateStatus(context.Context, uuid.UUID, order.Status) (order.Order, error)
	History(context.Context, uuid.UUID) ([]order.StatusChange, error)
}

type Handler struct {
	service OrderService
	logger  *slog.Logger
}

func NewRouter(service OrderService, logger *slog.Logger) http.Handler {
	h := &Handler{service: service, logger: logger}
	router := chi.NewRouter()
	router.Use(requestLoggingMiddleware(logger))
	router.Use(recoveryMiddleware(logger))
	router.Get("/health", h.health)
	router.Get("/swagger/*", httpSwagger.Handler())
	router.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", h.create)
		r.Get("/", h.list)
		r.Get("/{id}", h.get)
		r.Patch("/{id}/status", h.updateStatus)
		r.Get("/{id}/history", h.history)
	})
	return router
}

// health godoc
// @Summary Health check
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// create godoc
// @Summary Create order
// @Accept json
// @Produce json
// @Param request body createOrderRequest true "Order items"
// @Success 201 {object} orderResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/v1/orders [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var request createOrderRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid JSON", "INVALID_REQUEST"})
		return
	}
	items := make([]order.Item, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, order.Item{Name: item.Name, Quantity: item.Quantity, Price: item.Price})
	}
	value, err := h.service.Create(r.Context(), order.CreateOrderInput{Items: items})
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toOrderResponse(value))
}

// list godoc
// @Summary List orders
// @Produce json
// @Success 200 {array} orderResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/orders [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	values, err := h.service.List(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}
	response := make([]orderResponse, 0, len(values))
	for _, value := range values {
		response = append(response, toOrderResponse(value))
	}
	writeJSON(w, http.StatusOK, response)
}

// get godoc
// @Summary Get order
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} orderResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /api/v1/orders/{id} [get]
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	value, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toOrderResponse(value))
}

// updateStatus godoc
// @Summary Update order status
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body updateStatusRequest true "New status"
// @Success 200 {object} orderResponse
// @Failure 400,404,409,500 {object} errorResponse
// @Router /api/v1/orders/{id}/status [patch]
func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var request updateStatusRequest
	if err := decodeJSON(w, r, &request); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid JSON", "INVALID_REQUEST"})
		return
	}
	value, err := h.service.UpdateStatus(r.Context(), id, request.Status)
	if err != nil {
		h.writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toOrderResponse(value))
}

// history godoc
// @Summary Get order status history
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {array} statusChangeResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /api/v1/orders/{id}/history [get]
func (h *Handler) history(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	value, err := h.service.History(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}
	response := make([]statusChangeResponse, 0, len(value))
	for _, change := range value {
		response = append(response, statusChangeResponse{Status: change.Status, ChangedAt: change.ChangedAt})
	}
	writeJSON(w, http.StatusOK, response)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, value any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("expected a single JSON value")
	}
	return nil
}

func parseID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil || id == uuid.Nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid order id", "INVALID_REQUEST"})
		return uuid.Nil, false
	}
	return id, true
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, order.ErrInvalidOrder):
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid order", "INVALID_ORDER"})
	case errors.Is(err, order.ErrInvalidStatus):
		writeJSON(w, http.StatusBadRequest, errorResponse{"invalid status", "INVALID_STATUS"})
	case errors.Is(err, order.ErrOrderNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{"order not found", "ORDER_NOT_FOUND"})
	case errors.Is(err, order.ErrInvalidStatusTransition):
		writeJSON(w, http.StatusConflict, errorResponse{"status transition is not allowed", "INVALID_STATUS_TRANSITION"})
	default:
		h.logger.Error("request failed", slog.Any("error", err))
		writeJSON(w, http.StatusInternalServerError, errorResponse{"internal server error", "INTERNAL_ERROR"})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
