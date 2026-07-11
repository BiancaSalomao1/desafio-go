package routes

/*
struct Router

Responsabilidades:
- registrar todas as rotas da API.

Campos:
- mux

Métodos:
- NewRouter()
*/

import (
	"net/http"

	"desafio-go/infrastructure/http/handler"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	customerHandler *handler.CustomerHandler,
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /products", productHandler.Create)
	mux.HandleFunc("GET /products", productHandler.List)
	mux.HandleFunc("GET /products/{id}", productHandler.GetByID)

	mux.HandleFunc("POST /customers", customerHandler.Create)
	mux.HandleFunc("GET /customers", customerHandler.List)
	mux.HandleFunc("GET /customers/{id}", customerHandler.GetByID)

	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", userHandler.List)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)

	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetByID)
	mux.HandleFunc("PATCH /orders/{id}/pay", orderHandler.Pay)
	mux.HandleFunc("PATCH /orders/{id}/cancel", orderHandler.Cancel)

	return mux
}
