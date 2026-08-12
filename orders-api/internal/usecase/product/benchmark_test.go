package product

/*
Benchmark do CreateProductUseCase.

Responsabilidades:
- medir desempenho da criação de produtos.
*/

import (
	"testing"

	"desafio-go/orders-api/internal/domain"
)

func BenchmarkUpdateProductUseCase_Execute(b *testing.B) {

	repository := &productRepositoryStub{
		product: &domain.Product{
			ID:    "1",
			Name:  "Notebook",
			Price: 3500,
			Stock: 10,
		},
	}

	useCase := NewUpdateProductUseCase(repository)

	product := &domain.Product{
		ID:    "1",
		Name:  "Notebook Gamer",
		Price: 5000,
		Stock: 15,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(product)
	}
}
