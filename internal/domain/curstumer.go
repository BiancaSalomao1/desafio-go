package domain

/*
struct Customer

- identificar o cliente;
- armazenar nome;
- armazenar e-mail.

Métodos:
- construtor NewCustomer()
- Validate()
*/

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewCustomer(id, name, email string) *Customer {
	return &Customer{
		ID:    id,
		Name:  name,
		Email: email,
	}
}

func (c *Customer) Validate() error {
	if c.Name == "" {
		return ErrCustomerInvalid
	}

	if c.Email == "" {
		return ErrCustomerInvalid
	}

	return nil
}
