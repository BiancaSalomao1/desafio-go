package routes

/*
struct Router

Responsabilidades:
- registrar todas as rotas da API;
- definir quais rotas são públicas;
- proteger as demais com JWT.

Métodos:
- NewRouter()
*/

import (
	"net/http"

	_ "orders-api/docs"

	"orders-api/infrastructure/http/handler"
	"orders-api/infrastructure/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	productHandler *handler.ProductHandler,
	customerHandler *handler.CustomerHandler,
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
	authHandler *handler.AuthHandler,
	jwtSecret string,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Middleware JWT
	auth := middleware.Auth(jwtSecret)

	// ==========================
	// Rotas Públicas
	// ==========================

	mux.HandleFunc("POST /login", authHandler.Login)

	// Permite criar o primeiro usuário
	mux.HandleFunc("POST /users", userHandler.Create)

	// ==========================
	// Produtos
	// ==========================

	mux.Handle("POST /products",
		auth(http.HandlerFunc(productHandler.Create)))

	mux.Handle("GET /products",
		auth(http.HandlerFunc(productHandler.List)))

	mux.Handle("GET /products/{id}",
		auth(http.HandlerFunc(productHandler.GetByID)))

	mux.Handle("PUT /products/{id}",
		auth(http.HandlerFunc(productHandler.Update)))

	mux.Handle("DELETE /products/{id}",
		auth(http.HandlerFunc(productHandler.Delete)))

	// ==========================
	// Clientes
	// ==========================

	mux.Handle("POST /customers",
		auth(http.HandlerFunc(customerHandler.Create)))

	mux.Handle("GET /customers",
		auth(http.HandlerFunc(customerHandler.List)))

	mux.Handle("GET /customers/{id}",
		auth(http.HandlerFunc(customerHandler.GetByID)))

	mux.Handle("PUT /customers/{id}",
		auth(http.HandlerFunc(customerHandler.Update)))

	mux.Handle("DELETE /customers/{id}",
		auth(http.HandlerFunc(customerHandler.Delete)))

	// ==========================
	// Usuários
	// ==========================

	mux.Handle("GET /users",
		auth(http.HandlerFunc(userHandler.List)))

	mux.Handle("GET /users/{id}",
		auth(http.HandlerFunc(userHandler.GetByID)))

	mux.Handle("PUT /users/{id}",
		auth(http.HandlerFunc(userHandler.Update)))

	mux.Handle("DELETE /users/{id}",
		auth(http.HandlerFunc(userHandler.Delete)))

	// ==========================
	// Pedidos
	// ==========================

	mux.Handle("POST /orders",
		auth(http.HandlerFunc(orderHandler.Create)))

	mux.Handle("GET /orders",
		auth(http.HandlerFunc(orderHandler.List)))

	mux.Handle("GET /orders/{id}",
		auth(http.HandlerFunc(orderHandler.GetByID)))

	mux.Handle("PATCH /orders/{id}/pay",
		auth(http.HandlerFunc(orderHandler.Pay)))

	mux.Handle("PATCH /orders/{id}/cancel",
		auth(http.HandlerFunc(orderHandler.Cancel)))

	return mux
}
