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

	"desafio-go/infrastructure/http/httpx"

	authdto "desafio-go/internal/dto/auth"
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

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request authdto.LoginRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	token, err := h.loginUseCase.Execute(
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
		authdto.LoginResponse{
			AccessToken: token,
		},
	)
}
