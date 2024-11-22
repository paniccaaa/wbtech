package app

import (
	"log/slog"
	"net/http"

	api "github.com/paniccaaa/wbtech/internal/api/order"
)

func InitRouter(orderService api.GetProvider, log *slog.Logger) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /order/{order_uid}", api.HandleGetOrder(orderService, log))

	return r
}
