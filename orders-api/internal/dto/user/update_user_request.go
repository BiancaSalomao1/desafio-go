package user

/*
struct UpdateUserRequest

Responsabilidades:
- receber os dados para atualização de um usuário.

Campos:
- name
- email

*/

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
