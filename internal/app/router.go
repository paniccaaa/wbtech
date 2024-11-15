package app

import (
	"net/http"

	api "github.com/paniccaaa/wbtech/internal/api/order"
)

func InitRouter(orderService api.OrderService) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /order", api.HandleGetOrder(orderService))

	return r
}
