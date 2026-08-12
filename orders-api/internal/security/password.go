package security

/*
Função HashPassword

Responsabilidades:
- gerar o hash de uma senha utilizando bcrypt.

Métodos:
- HashPassword()
*/

import "golang.org/x/crypto/bcrypt"

func HashPassword(
	password string,
) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

/*
Função CheckPassword

Responsabilidades:
- comparar uma senha em texto puro com um hash bcrypt.

Métodos:
- CheckPassword()
*/

func CheckPassword(
	hash string,
	password string,
) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
