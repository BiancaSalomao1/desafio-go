package mapper

import (
	"desafio-go/orders-api/internal/domain"
	productdto "desafio-go/orders-api/internal/dto/product"
	"testing"
)

func TestToProduct(t *testing.T) {

	request := productdto.CreateProductRequest{
		Name:  "Notebook",
		Price: 4999.90,
		Stock: 10,
	}

	product := ToProduct(request)

	if product.ID == "" {
		t.Fatal("expected generated id")
	}

	if product.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, product.Name)
	}

	if product.Price != request.Price {
		t.Fatalf("expected %.2f, got %.2f", request.Price, product.Price)
	}

	if product.Stock != request.Stock {
		t.Fatalf("expected %d, got %d", request.Stock, product.Stock)
	}
}

func TestToProductUpdate(t *testing.T) {

	request := productdto.UpdateProductRequest{
		Name:  "Mouse",
		Price: 99.90,
		Stock: 5,
	}

	product := ToProductUpdate(
		"PROD001",
		request,
	)

	if product.ID != "PROD001" {
		t.Fatalf("expected PROD001, got %s", product.ID)
	}

	if product.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, product.Name)
	}

	if product.Price != request.Price {
		t.Fatalf("expected %.2f, got %.2f", request.Price, product.Price)
	}

	if product.Stock != request.Stock {
		t.Fatalf("expected %d, got %d", request.Stock, product.Stock)
	}
}

func TestToProductResponse(t *testing.T) {

	product := domain.NewProduct(
		"PROD001",
		"Notebook",
		4999.90,
		10,
	)

	response := ToProductResponse(product)

	if response.ID != product.ID {
		t.Fatal("id mismatch")
	}

	if response.Name != product.Name {
		t.Fatal("name mismatch")
	}

	if response.Price != product.Price {
		t.Fatal("price mismatch")
	}

	if response.Stock != product.Stock {
		t.Fatal("stock mismatch")
	}
}

func TestToProductResponseList(t *testing.T) {

	products := []*domain.Product{
		domain.NewProduct(
			"PROD001",
			"Notebook",
			4999.90,
			10,
		),
		domain.NewProduct(
			"PROD002",
			"Mouse",
			99.90,
			20,
		),
	}

	response := ToProductResponseList(products)

	if len(response) != 2 {
		t.Fatalf("expected 2 products, got %d", len(response))
	}

	if response[0].ID != "PROD001" {
		t.Fatal("unexpected first product")
	}

	if response[1].ID != "PROD002" {
		t.Fatal("unexpected second product")
	}
}

func TestToProductResponseList_Empty(t *testing.T) {

	response := ToProductResponseList([]*domain.Product{})

	if len(response) != 0 {
		t.Fatalf("expected empty slice, got %d", len(response))
	}
}
