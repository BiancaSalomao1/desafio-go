package customer

/*
struct CreateCustomerRequest

Responsabilidades:

- receber os dados para criação de um cliente.

Campos:
- name
- email
*/

type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
