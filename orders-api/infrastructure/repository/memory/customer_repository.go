package memory

/*
struct MemoryCustomerRepository

Responsabilidades:

- armazenar clientes em memória;
- implementar CustomerRepository.

Métodos:
- NewMemoryCustomerRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindByEmail()
- FindAll()
*/

import (
	"orders-api/internal/domain"
	"orders-api/internal/repository"
)

type MemoryCustomerRepository struct {
	customers map[string]*domain.Customer
}

func NewMemoryCustomerRepository() repository.CustomerRepository {
	return &MemoryCustomerRepository{
		customers: make(map[string]*domain.Customer),
	}
}

func (r *MemoryCustomerRepository) Save(customer *domain.Customer) error {

	if err := customer.Validate(); err != nil {
		return err
	}

	r.customers[customer.ID] = customer

	return nil
}

func (r *MemoryCustomerRepository) Update(customer *domain.Customer) error {

	if _, exists := r.customers[customer.ID]; !exists {
		return domain.ErrCustomerNotFound
	}

	if err := customer.Validate(); err != nil {
		return err
	}

	r.customers[customer.ID] = customer

	return nil
}

func (r *MemoryCustomerRepository) Delete(id string) error {

	if _, exists := r.customers[id]; !exists {
		return domain.ErrCustomerNotFound
	}

	delete(r.customers, id)

	return nil
}

func (r *MemoryCustomerRepository) FindByID(id string) (*domain.Customer, error) {

	customer, exists := r.customers[id]

	if !exists {
		return nil, domain.ErrCustomerNotFound
	}

	return customer, nil
}

func (r *MemoryCustomerRepository) FindByEmail(email string) (*domain.Customer, error) {

	for _, customer := range r.customers {

		if customer.Email == email {
			return customer, nil
		}
	}

	return nil, domain.ErrCustomerNotFound
}

func (r *MemoryCustomerRepository) FindAll() ([]*domain.Customer, error) {

	customers := make([]*domain.Customer, 0, len(r.customers))

	for _, customer := range r.customers {
		customers = append(customers, customer)
	}

	return customers, nil
}
