package user

/*
struct UserResponse

Responsabilidades:

- retornar os dados de um usuário.

Campos:
- id
- name
- email
*/

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
