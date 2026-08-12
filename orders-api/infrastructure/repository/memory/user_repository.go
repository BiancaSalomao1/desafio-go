package memory

/*
struct MemoryUserRepository

Responsabilidades:

- armazenar usuários em memória;
- implementar UserRepository.

Métodos:
- NewMemoryUserRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindByEmail()
- FindAll()
*/

import (
	"desafio-go/orders-api/internal/domain"
	"desafio-go/orders-api/internal/repository"
)

type MemoryUserRepository struct {
	users map[string]*domain.User
}

func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *MemoryUserRepository) Save(user *domain.User) error {

	if err := user.Validate(); err != nil {
		return err
	}

	r.users[user.ID] = user

	return nil
}

func (r *MemoryUserRepository) Update(user *domain.User) error {

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	if err := user.Validate(); err != nil {
		return err
	}

	r.users[user.ID] = user

	return nil
}

func (r *MemoryUserRepository) Delete(id string) error {

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *MemoryUserRepository) FindByID(id string) (*domain.User, error) {

	user, exists := r.users[id]

	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *MemoryUserRepository) FindByEmail(email string) (*domain.User, error) {

	for _, user := range r.users {

		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func (r *MemoryUserRepository) FindAll() ([]*domain.User, error) {

	users := make([]*domain.User, 0, len(r.users))

	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}
