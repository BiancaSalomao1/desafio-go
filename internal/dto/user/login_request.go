package user

/*
struct LoginRequest

Responsabilidades:

- receber os dados para autenticação.

Campos:
- email
- password
*/

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
