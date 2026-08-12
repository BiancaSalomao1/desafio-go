package auth

/*
struct LoginResponse

Responsabilidades:
- retornar o token JWT.

Campos:
- accessToken
*/

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}
