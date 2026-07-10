package customer

/*
struct UpdateCustomerRequest

Responsabilidades:
- receber os dados para atualização de um cliente.

Campos:
- name
- email

*/

type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
