package customer

/*
Benchmark do CreateCustomerUseCase.

Responsabilidades:
- medir desempenho da criação de clientes.
*/

import (
	"testing"

	"desafio-go/internal/domain"
)

func BenchmarkCreateCustomerUseCase_Execute(b *testing.B) {

	repository := &customerRepositoryStub{}

	useCase := NewCreateCustomerUseCase(repository)

	customer := &domain.Customer{
		Name: "João da Silva",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(customer)
	}
}
