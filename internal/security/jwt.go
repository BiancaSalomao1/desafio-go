package security

/*
Função GenerateToken

Responsabilidades:
- gerar JWT.

Métodos:
- GenerateToken()
*/

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(
	userID string,
	email string,
	secret string,
	expiration time.Duration,
) (string, error) {

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(expiration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(secret),
	)
}
