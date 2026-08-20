package handler

/*
struct AuthHandler

Responsabilidades:
- autenticar usuário.

Métodos:
- Login()
*/

import (
	"net/http"

	"orders-api/infrastructure/http/httpx"

	"orders-api/internal/dto/auth"
)

type AuthHandler struct {
	loginUseCase LoginUseCase
}

func NewAuthHandler(
	loginUseCase LoginUseCase,
) *AuthHandler {

	return &AuthHandler{
		loginUseCase: loginUseCase,
	}
}

// Login
//
// @Summary Login
// @Description Autentica um usuário e retorna um JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Credenciais"
// @Success 200 {object} auth.LoginResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request auth.LoginRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	token, err := h.loginUseCase.Execute(
		r.Context(),
		request.Email,
		request.Password,
	)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		auth.LoginResponse{
			AccessToken: token,
		},
	)
}
