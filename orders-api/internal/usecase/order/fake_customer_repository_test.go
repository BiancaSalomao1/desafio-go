package order

/*
FakeCustomerRepository.

Responsabilidades:
- simular o CustomerRepository;
- controlar retornos;
- registrar chamadas.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

import (
	"orders-api/internal/domain"
)

type FakeCustomerRepository struct {
	customer *domain.Customer

	findError error

	findCalls int

	spy  *OrderSpy
	mock *OrderMock
}

func (r *FakeCustomerRepository) Save(customer *domain.Customer) error {
	return nil
}

func (r *FakeCustomerRepository) Update(customer *domain.Customer) error {
	return nil
}

func (r *FakeCustomerRepository) Delete(id string) error {
	return nil
}

func (r *FakeCustomerRepository) FindByID(id string) (*domain.Customer, error) {

	r.findCalls++

	if r.spy != nil {
		r.spy.Add("customer.find")
	}

	if r.mock != nil {
		r.mock.CustomerFind++
	}

	if r.findError != nil {
		return nil, r.findError
	}

	return r.customer, nil
}

func (r *FakeCustomerRepository) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}
func (r *FakeCustomerRepository) FindByEmail(
	email string,
) (*domain.Customer, error) {

	if r.customer != nil && r.customer.Email == email {
		return r.customer, nil
	}

	return nil, domain.ErrCustomerNotFound
}
