package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/paniccaaa/wbtech/internal/model"
)

type GetProvider interface {
	GetOrder(ctx context.Context, orderUID model.OrderUID) (model.Order, error)
}

func HandleGetOrder(orderService GetProvider, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.PathValue("order_uid")
		if orderUID == "" {
			http.Error(w, "order_uid is required", http.StatusBadRequest)
			return
		}

		uid := model.OrderUID(orderUID)

		log.Info("get orderUID from path", slog.String("uid", string(uid)))

		order, err := orderService.GetOrder(r.Context(), uid)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				log.Warn("order UID not found", slog.String("orderUID", orderUID))

				http.Error(w, "order not found", http.StatusNotFound)
				return
			}

			log.Error("failed to get order", slog.String("err", err.Error()))

			http.Error(w, fmt.Sprintf("failed to get order: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Error("failed to encode order", slog.String("orderUID", orderUID))

			http.Error(w, "failed to encode order", http.StatusInternalServerError)
			return
		}
	}
}
