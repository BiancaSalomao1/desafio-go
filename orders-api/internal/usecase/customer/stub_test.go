package customer

/*
Stub do CustomerRepository.

Responsabilidades:
- simular o repositório;
- controlar o retorno do método Save().
*/

import (
	"orders-api/internal/domain"
)

type customerRepositoryStub struct {
	customer *domain.Customer

	saveError   error
	findError   error
	updateError error
}

func (r *customerRepositoryStub) Save(customer *domain.Customer) error {
	return r.saveError
}

func (r *customerRepositoryStub) Update(customer *domain.Customer) error {
	return r.updateError
}

func (r *customerRepositoryStub) Delete(id string) error {
	return nil
}

func (r *customerRepositoryStub) FindByID(id string) (*domain.Customer, error) {

	if r.findError != nil {
		return nil, r.findError
	}

	return r.customer, nil
}

func (r *customerRepositoryStub) FindAll() ([]*domain.Customer, error) {
	return nil, nil
}
func (r *customerRepositoryStub) FindByEmail(
	email string,
) (*domain.Customer, error) {

	if r.customer != nil && r.customer.Email == email {
		return r.customer, nil
	}

	return nil, domain.ErrCustomerNotFound
}
