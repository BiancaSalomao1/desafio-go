package repository

import "desafio-go/internal/domain"

/*
interface ProductRepository

Responsabilidades:

- salvar produto;
- atualizar produto;
- remover produto;
- buscar produto por ID;
- listar produtos.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

type ProductRepository interface {
	Save(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
	FindByID(id string) (*domain.Product, error)
	FindAll() ([]*domain.Product, error)
}
