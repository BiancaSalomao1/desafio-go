package app

import (
	"net/http"

	"product-service/config"
	"product-service/infrastructure/database"
	"product-service/infrastructure/http/handler"
	"product-service/infrastructure/http/routes"
	"product-service/infrastructure/repository/postgres"

	productusecase "product-service/internal/usecase/product"
)

// NewApplication monta toda a aplicação do Product Service.
//
// Fluxo:
//
// HTTP
// ↓
// Router
// ↓
// Handler
// ↓
// Use Cases
// ↓
// ProductRepository
// ↓
// PostgreSQL
func NewApplication(
	_ *config.Config,
	db *database.Database,
) (http.Handler, error) {

	productRepository := postgres.NewProductRepository(
		db.Pool,
	)

	createProduct := productusecase.NewCreateProductUseCase(
		productRepository,
	)

	getProduct := productusecase.NewGetProductUseCase(
		productRepository,
	)

	listProducts := productusecase.NewListProductsUseCase(
		productRepository,
	)

	updateProduct := productusecase.NewUpdateProductUseCase(
		productRepository,
	)

	deleteProduct := productusecase.NewDeleteProductUseCase(
		productRepository,
	)

	productHandler := handler.NewProductHandler(
		createProduct,
		getProduct,
		listProducts,
		updateProduct,
		deleteProduct,
	)

	return routes.NewRouter(productHandler), nil
}
