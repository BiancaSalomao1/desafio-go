package mapper

import (
	"desafio-go/internal/domain"
	productdto "desafio-go/internal/dto/product"

	"github.com/google/uuid"
)

func ToProduct(request productdto.CreateProductRequest) *domain.Product {

	return domain.NewProduct(
		uuid.NewString(),
		request.Name,
		request.Price,
		request.Stock,
	)
}

func ToUpdatedProduct(
	id string,
	request productdto.UpdateProductRequest,
) *domain.Product {

	return domain.NewProduct(
		id,
		request.Name,
		request.Price,
		request.Stock,
	)
}

func ToProductResponse(
	product *domain.Product,
) productdto.ProductResponse {

	return productdto.ProductResponse{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}
}

func ToProductResponseList(
	products []*domain.Product,
) []productdto.ProductResponse {

	response := make(
		[]productdto.ProductResponse,
		0,
		len(products),
	)

	for _, product := range products {
		response = append(
			response,
			ToProductResponse(product),
		)
	}

	return response
}
