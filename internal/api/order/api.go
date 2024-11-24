package order

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/paniccaaa/wbtech/internal/model"
)

//go:generate mockery --name GetProvider
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

		orderTemplate := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Order Details</title>
			</head>
			<body>
				<h1>Order Details</h1>
				<p><strong>Order ID:</strong> {{.OrderUID}}</p>
				<p><strong>Track Number:</strong> {{.TrackNumber}}</p>
				<p><strong>Customer ID:</strong> {{.CustomerID}}</p>
				<p><strong>Status:</strong> {{.Payment.Transaction}}</p>
				<p><strong>Total Amount:</strong> ${{.Payment.Amount}}</p>
				<p><strong>Items:</strong></p>
				<ul>
					{{range .Items}}
						<li>{{.Name}} ({{.Price}}) - Quantity: {{.Sale}} - Total Price: {{.TotalPrice}}</li>
					{{end}}
				</ul>
			</body>
			</html>
		`

		// Создаем шаблон
		tmpl, err := template.New("order").Parse(orderTemplate)
		if err != nil {
			log.Error("failed to parse template", slog.String("err", err.Error()))
			http.Error(w, fmt.Sprintf("failed to parse template: %v", err), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, order); err != nil {
			log.Error("failed to render template", slog.String("err", err.Error()))
			http.Error(w, fmt.Sprintf("failed to render template: %v", err), http.StatusInternalServerError)
		}

		// w.Header().Set("Content-Type", "application/json")
		// if err := json.NewEncoder(w).Encode(order); err != nil {
		// 	log.Error("failed to encode order", slog.String("orderUID", orderUID), slog.String("err", err.Error()))

		// 	http.Error(w, "failed to encode order", http.StatusInternalServerError)
		// 	return
		// }
	}
}
