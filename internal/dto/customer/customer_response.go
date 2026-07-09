package customer

/*
struct CustomerResponse

Responsabilidades:

- retornar os dados de um cliente.

Campos:
- id
- name
- email
*/

type CustomerResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
