package repository

import "orders-api/internal/domain"

/*
interface UserRepository

Responsabilidades:

- salvar usuário;
- atualizar usuário;
- remover usuário;
- buscar usuário por ID;
- buscar usuário por e-mail;
- listar usuários.

Métodos:
- Save()
- Update()
- Delete()
- FindByID()
- FindByEmail()
- FindAll()
*/

type UserRepository interface {
	Save(user *domain.User) error
	Update(user *domain.User) error
	Delete(id string) error
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]*domain.User, error)
}
