package order

/*
Mock do CreateOrderUseCase.

Responsabilidades:
- validar quantidade de chamadas dos repositórios.
*/

import "errors"

type OrderMock struct {
	ExpectedCustomerFind  int
	ExpectedProductFind   int
	ExpectedProductUpdate int
	ExpectedOrderSave     int

	CustomerFind  int
	ProductFind   int
	ProductUpdate int
	OrderSave     int
}

func (m *OrderMock) Verify() error {

	if m.CustomerFind != m.ExpectedCustomerFind {
		return errors.New("unexpected customer repository calls")
	}

	if m.ProductFind != m.ExpectedProductFind {
		return errors.New("unexpected product find calls")
	}

	if m.ProductUpdate != m.ExpectedProductUpdate {
		return errors.New("unexpected product update calls")
	}

	if m.OrderSave != m.ExpectedOrderSave {
		return errors.New("unexpected order save calls")
	}

	return nil
}
