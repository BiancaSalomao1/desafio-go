package domain

/*
struct User

- identificar o usuário;
- armazenar nome;
- armazenar e-mail;
- armazenar passwordHash.

Métodos:
- construtor NewUser()
- Validate()
- CheckPassword()
- Update()
*/

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func NewUser(id, name, email, passwordHash string) *User {
	return &User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrUserNameRequired
	}

	if u.Email == "" {
		return ErrUserEmailRequired
	}

	if u.PasswordHash == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (u *User) CheckPassword(password string) bool {
	return u.PasswordHash == password
}

func (u *User) Update(
	name string,
	email string,
) error {

	u.Name = name
	u.Email = email

	return u.Validate()
}
