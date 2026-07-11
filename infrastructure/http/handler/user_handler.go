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
- updateUserUseCase
- deleteUserUseCase

Métodos:
- NewUserHandler()
- Create()
- GetByID()
- List()
- Update()
- Delete()
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
	updateUserUseCase UpdateUserUseCase
	deleteUserUseCase DeleteUserUseCase
}

func NewUserHandler(
	createUserUseCase CreateUserUseCase,
	getUserUseCase GetUserUseCase,
	listUsersUseCase ListUsersUseCase,
	updateUserUseCase UpdateUserUseCase,
	deleteUserUseCase DeleteUserUseCase,
) *UserHandler {

	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
		listUsersUseCase:  listUsersUseCase,
		updateUserUseCase: updateUserUseCase,
		deleteUserUseCase: deleteUserUseCase,
	}
}

// Create
//
// @Summary Criar Usuário
// @Description Cadastra um novo usuário.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body user.CreateUserRequest true "Usuário"
// @Success 201 {object} user.UserResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Router /users [post]
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

// GetByID
//
// @Summary Buscar Usuário
// @Description Busca um usuário pelo ID.
// @Tags Users
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} user.UserResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [get]
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

// List
//
// @Summary Listar Usuários
// @Description Retorna todos os usuários.
// @Tags Users
// @Produce json
// @Success 200 {array} user.UserResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /users [get]
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

// Update
//
// @Summary Atualizar Usuário
// @Description Atualiza um usuário.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body user.UpdateUserRequest true "Usuário"
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	var request userdto.UpdateUserRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	user := mapper.ToUserUpdate(
		id,
		request,
	)

	if err := h.updateUserUseCase.Execute(user); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
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

// Delete
//
// @Summary Remover Usuário
// @Description Remove um usuário.
// @Tags Users
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	if err := h.deleteUserUseCase.Execute(id); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	httpx.WriteStatus(
		w,
		http.StatusNoContent,
	)
}
