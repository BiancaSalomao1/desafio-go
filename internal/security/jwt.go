package security

/*
Função GenerateToken

Responsabilidades:
- gerar JWT.
- validar um JWT.
- retornar os claims.

Métodos:
- GenerateToken()
- ValidateToken()
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

func ValidateToken(
	tokenString string,
	secret string,
) (jwt.MapClaims, error) {

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
