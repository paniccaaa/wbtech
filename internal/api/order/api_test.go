package order

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paniccaaa/wbtech/internal/api/order/mocks"
	"github.com/paniccaaa/wbtech/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGetOrder(t *testing.T) {
	mockOrderService := new(mocks.GetProvider)
	logger := slog.Default()

	tests := []struct {
		name               string
		orderUID           string
		mockGetOrderReturn func() (model.Order, error)
		wantStatusCode     int
		wantBody           string
	}{
		{
			name:     "order not found",
			orderUID: "12345",
			mockGetOrderReturn: func() (model.Order, error) {
				return model.Order{}, model.ErrOrderNotFound
			},
			wantStatusCode: http.StatusNotFound,
			wantBody:       "order not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.orderUID != "" {
				mockOrderService.On("GetOrder", mock.Anything, model.OrderUID(tt.orderUID)).
					Return(tt.mockGetOrderReturn())
			}

			req, err := http.NewRequest(http.MethodGet, "/orders/"+tt.orderUID, nil)
			assert.NoError(t, err)

			r := http.NewServeMux()
			r.HandleFunc("/orders/{order_uid}", HandleGetOrder(mockOrderService, logger))

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatusCode, rr.Code)
			assert.Equal(t, tt.wantBody, rr.Body.String())
		})
	}
}
