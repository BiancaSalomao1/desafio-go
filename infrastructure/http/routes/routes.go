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
	"desafio-go/infrastructure/http/handler"
	"net/http"

	_ "desafio-go/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	customerHandler *handler.CustomerHandler,
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
	authHandler *handler.AuthHandler,
) *http.ServeMux {

	mux := http.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("POST /products", productHandler.Create)
	mux.HandleFunc("GET /products", productHandler.List)
	mux.HandleFunc("GET /products/{id}", productHandler.GetByID)
	mux.HandleFunc("PUT /products/{id}", productHandler.Update)
	mux.HandleFunc("DELETE /products/{id}", productHandler.Delete)

	mux.HandleFunc("POST /customers", customerHandler.Create)
	mux.HandleFunc("GET /customers", customerHandler.List)
	mux.HandleFunc("GET /customers/{id}", customerHandler.GetByID)
	mux.HandleFunc("PUT /customers/{id}", customerHandler.Update)
	mux.HandleFunc("DELETE /custumers/{id}", customerHandler.Delete)

	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", userHandler.List)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders", orderHandler.List)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetByID)
	mux.HandleFunc("PATCH /orders/{id}/pay", orderHandler.Pay)
	mux.HandleFunc("PATCH /orders/{id}/cancel", orderHandler.Cancel)

	mux.HandleFunc("POST /login", authHandler.Login)

	return mux
}
