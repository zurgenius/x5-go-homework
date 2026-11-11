package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/order-service/internal/handler"
	"example.com/order-service/internal/repository/memory"
	"example.com/order-service/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOrderAPI(t *testing.T) {
	router := handler.NewRouter(service.NewOrderService(memory.NewOrderRepository()), slog.New(slog.NewTextHandler(io.Discard, nil)))
	request := func(method, path, body string, status int) map[string]any {
		t.Helper()
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(method, path, bytes.NewBufferString(body)))
		require.Equal(t, status, response.Code, response.Body.String())
		var result map[string]any
		if len(response.Body.Bytes()) > 0 && response.Body.Bytes()[0] == '{' {
			require.NoError(t, json.Unmarshal(response.Body.Bytes(), &result))
		}
		return result
	}
	request(http.MethodGet, "/health", "", 200)
	created := request(http.MethodPost, "/api/v1/orders", `{"items":[{"name":"item","quantity":2,"price":10}]}`, 201)
	require.EqualValues(t, 20, created["total_price"])
	id := created["id"].(string)
	parsed, err := uuid.Parse(id)
	require.NoError(t, err)
	require.Equal(t, uuid.Version(7), parsed.Version())
	request(http.MethodGet, "/api/v1/orders", "", 200)
	request(http.MethodGet, "/api/v1/orders/"+id, "", 200)
	require.Equal(t, "INVALID_STATUS_TRANSITION", request(http.MethodPatch, "/api/v1/orders/"+id+"/status", `{"status":"completed"}`, 409)["code"])
	request(http.MethodPatch, "/api/v1/orders/"+id+"/status", `{"status":"processing"}`, 200)
	request(http.MethodGet, "/api/v1/orders/"+id+"/history", "", 200)
	request(http.MethodGet, "/api/v1/orders/"+uuid.NewString(), "", 404)
	request(http.MethodGet, "/api/v1/orders/not-a-uuid", "", 400)
	request(http.MethodPost, "/api/v1/orders", `{}`, 400)
	request(http.MethodPost, "/api/v1/orders", `{} {}`, 400)
}

func TestRecoveryAndSwagger(t *testing.T) {
	var logs bytes.Buffer
	router := handler.NewRouter(nil, slog.New(slog.NewJSONHandler(&logs, nil)))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil))
	require.Equal(t, http.StatusInternalServerError, response.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	require.Equal(t, "INTERNAL_ERROR", body["code"])
	require.Contains(t, logs.String(), `"status":500`)
	require.Contains(t, logs.String(), `"request_id":`)

	for _, path := range []string{"/swagger/index.html", "/swagger/doc.json"} {
		response = httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		require.Equal(t, http.StatusOK, response.Code, path)
	}
}
