package order

import (
	"context"
	"net/http"

	"github.com/paniccaaa/wbtech/internal/model"
)

type OrderService interface {
	GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error)
}

func HandleGetOrder(orderService OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
	}
}
