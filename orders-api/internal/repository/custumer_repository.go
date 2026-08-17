package repository

import "orders-api/internal/domain"

/*
interface CustomerRepository

Responsabilidades:

- salvar cliente;
- atualizar cliente;
- remover cliente;
- buscar cliente por ID;
- listar clientes.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

type CustomerRepository interface {
	Save(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(id string) error
	FindByID(id string) (*domain.Customer, error)
	FindByEmail(email string) (*domain.Customer, error) // <- existe?
	FindAll() ([]*domain.Customer, error)
}
