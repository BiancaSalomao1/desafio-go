package handler

/*
struct UserHandler

Responsabilidades:
- receber requisições HTTP relacionadas a usuários;
- converter DTOs para entidades de domínio;
- executar os casos de uso;
- retornar respostas HTTP em JSON.

Campos:
- createUserUseCase
- getUserUseCase
- listUsersUseCase

Métodos:
- NewUserHandler()
- Create()
- GetByID()
- List()
*/

import (
	"net/http"

	"desafio-go/infrastructure/http/httpx"

	userdto "desafio-go/internal/dto/user"
	"desafio-go/internal/mapper"
)

type UserHandler struct {
	createUserUseCase CreateUserUseCase
	getUserUseCase    GetUserUseCase
	listUsersUseCase  ListUsersUseCase
}

func NewUserHandler(
	createUserUseCase CreateUserUseCase,
	getUserUseCase GetUserUseCase,
	listUsersUseCase ListUsersUseCase,
) *UserHandler {

	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
		listUsersUseCase:  listUsersUseCase,
	}
}

func (h *UserHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request userdto.CreateUserRequest

	if err := httpx.ReadJSON(r, &request); err != nil {
		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	user := mapper.ToUser(request)

	if err := h.createUserUseCase.Execute(user); err != nil {
		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)
		return
	}

	httpx.WriteJSON(
		w,
		http.StatusCreated,
		mapper.ToUserResponse(user),
	)
}

func (h *UserHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	user, err := h.getUserUseCase.Execute(id)
	if err != nil {
		httpx.WriteError(
			w,
			http.StatusNotFound,
			err,
		)
		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		mapper.ToUserResponse(user),
	)
}

func (h *UserHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	users, err := h.listUsersUseCase.Execute()
	if err != nil {
		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		mapper.ToUserResponseList(users),
	)
}
