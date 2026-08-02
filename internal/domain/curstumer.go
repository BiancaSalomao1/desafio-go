package domain

/*
struct Customer

Responsabilidades:
- identificar o cliente;
- armazenar nome;
- armazenar e-mail;
- armazenar o hash da senha;
- validar os dados;
- verificar a senha.

Métodos:
- NewCustomer()
- Validate()
- Update()
- CheckPassword()
*/

import (
	"desafio-go/internal/security"
)

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	Password string `json:"-"`

	PasswordHash string `json:"password_hash"`
}

func NewCustomer(
	id,
	name,
	email,
	passwordHash string,
) *Customer {

	return &Customer{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
}

func (c *Customer) Validate() error {

	if c.Name == "" {
		return ErrCustomerInvalid
	}

	if c.Email == "" {
		return ErrCustomerInvalid
	}

	if c.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (c *Customer) Update(
	name,
	email string,
) error {

	c.Name = name
	c.Email = email

	return c.Validate()
}

func (c *Customer) CheckPassword(
	password string,
) bool {

	return security.CheckPassword(
		c.PasswordHash,
		password,
	)
}
