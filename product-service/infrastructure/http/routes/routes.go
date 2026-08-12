package routes

import (
	"net/http"

	"product-service/infrastructure/http/handler"
)

func NewRouter(
	productHandler *handler.ProductHandler,
) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /products",
		productHandler.Create,
	)

	mux.HandleFunc(
		"GET /products",
		productHandler.List,
	)

	mux.HandleFunc(
		"GET /products/{id}",
		productHandler.GetByID,
	)

	mux.HandleFunc(
		"PUT /products/{id}",
		productHandler.Update,
	)

	mux.HandleFunc(
		"DELETE /products/{id}",
		productHandler.Delete,
	)

	return mux
}
