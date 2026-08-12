package auth

/*
struct LoginRequest

Responsabilidades:
- receber credenciais do usuário.

Campos:
- email
- password
*/

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
