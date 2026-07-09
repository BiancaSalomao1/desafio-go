package user

/*
struct CreateUserRequest

Responsabilidades:

- receber os dados para criação de um usuário.

Campos:
- name
- email
- password
*/

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
