package product

/*
Benchmark do CreateProductUseCase.

Responsabilidades:
- medir desempenho da criação de produtos.
*/

import (
	"testing"

	"desafio-go/internal/domain"
)

func BenchmarkCreateProductUseCase_Execute(b *testing.B) {

	repository := &productRepositoryStub{}

	useCase := NewCreateProductUseCase(repository)

	product := &domain.Product{
		Name:  "Notebook",
		Price: 3500,
		Stock: 10,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(product)
	}
}
