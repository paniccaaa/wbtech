package app

import (
	"net/http"

	api "github.com/paniccaaa/wbtech/internal/api/order"
)

func InitRouter(orderService api.GetProvider) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /order/{order_uid}", api.HandleGetOrder(orderService))

	return r
}
