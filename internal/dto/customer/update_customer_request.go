package customer

/*
struct UpdateCustomerRequest

Responsabilidades:

- receber os dados para atualização de um cliente.

Campos:
- id
- name
- email
*/

type UpdateCustomerRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
