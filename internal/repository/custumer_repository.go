package repository

import "desafio-go/internal/domain"

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
	FindAll() ([]*domain.Customer, error)
}
