package handler

/*
struct AuthHandler

Responsabilidades:
- autenticar usuário.

Métodos:
- Login()
- Logout()
*/

import (
	"errors"
	"net/http"
	"strings"

	"orders-api/infrastructure/http/httpx"

	"orders-api/internal/dto/auth"
)

type AuthHandler struct {
	loginUseCase  LoginUseCase
	logoutUseCase LogoutUseCase
}

func NewAuthHandler(
	loginUseCase LoginUseCase,
	logoutUseCase LogoutUseCase,

) *AuthHandler {

	return &AuthHandler{
		loginUseCase:  loginUseCase,
		logoutUseCase: logoutUseCase,
	}
}

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

func (h *AuthHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {

	header := r.Header.Get("Authorization")

	if !strings.HasPrefix(
		header,
		"Bearer ",
	) {
		httpx.WriteError(
			w,
			http.StatusUnauthorized,
			errors.New("invalid authorization header"),
		)
		return
	}

	token := strings.TrimPrefix(
		header,
		"Bearer ",
	)

	if err := h.logoutUseCase.Execute(
		r.Context(),
		token,
	); err != nil {

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
