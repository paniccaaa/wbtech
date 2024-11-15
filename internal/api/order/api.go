package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/paniccaaa/wbtech/internal/model"
)

type GetProvider interface {
	GetOrder(ctx context.Context, order_uid model.OrderUID) (model.Order, error)
}

func HandleGetOrder(orderService GetProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.PathValue("order_uid")
		if orderUID == "" {
			http.Error(w, "order_uid is required", http.StatusBadRequest)
			return
		}

		uid := model.OrderUID(orderUID)

		order, err := orderService.GetOrder(r.Context(), uid)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				http.Error(w, "order not found", http.StatusNotFound)
				return
			}

			http.Error(w, fmt.Sprintf("failed to get order: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "failed to encode order", http.StatusInternalServerError)
			return
		}
	}
}
