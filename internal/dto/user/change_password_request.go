package user

/*
struct ChangePasswordRequest

Responsabilidades:
- receber os dados para alteração de senha.

Campos:
- currentPassword
- newPassword

*/

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
