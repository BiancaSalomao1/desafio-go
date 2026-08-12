package order

/*
Benchmark do CreateOrderUseCase.

Responsabilidades:
- medir desempenho da criação de pedidos.
*/

import "testing"

func BenchmarkCreateOrderUseCase_Execute(b *testing.B) {

	product := newProduct()
	customer := newCustomer()

	factory := newFactory(product, customer)

	useCase := newUseCase(factory)

	order := newOrder(
		customer.ID,
		product.ID,
		1,
	)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = useCase.Execute(order)
	}
}
