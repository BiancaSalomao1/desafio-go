package order

/*
Spy do CreateOrderUseCase.

Responsabilidades:
- registrar chamadas aos repositórios;
- validar sequência de execução;
- contabilizar chamadas.
*/

type OrderSpy struct {
	Calls []string
}

func (s *OrderSpy) Add(call string) {
	s.Calls = append(s.Calls, call)
}

func (s *OrderSpy) Equals(expected []string) bool {

	if len(s.Calls) != len(expected) {
		return false
	}

	for i := range expected {
		if s.Calls[i] != expected[i] {
			return false
		}
	}

	return true
}
