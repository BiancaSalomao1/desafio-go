package handler

/*
struct CustomerHandler

Responsabilidades:
- receber requisições HTTP relacionadas a clientes;
- converter DTOs para entidades de domínio;
- executar os casos de uso;
- retornar respostas HTTP em JSON.

Campos:
- createCustomerUseCase
- getCustomerUseCase
- listCustomersUseCase
- updateCustomerUseCase
- deleteCustomerUseCase

Métodos:
- NewCustomerHandler()
- Create()
- GetByID()
- List()
- Update()
*/

import (
	"net/http"

	"orders-api/infrastructure/http/httpx"

	customerdto "orders-api/internal/dto/customer"
	"orders-api/internal/mapper"
)

type CustomerHandler struct {
	createCustomerUseCase CreateCustomerUseCase
	getCustomerUseCase    GetCustomerUseCase
	listCustomersUseCase  ListCustomersUseCase
	updateCustomerUseCase UpdateCustomerUseCase
	deleteCustomerUseCase DeleteCustomerUseCase
}

func NewCustomerHandler(
	createCustomerUseCase CreateCustomerUseCase,
	getCustomerUseCase GetCustomerUseCase,
	listCustomersUseCase ListCustomersUseCase,
	updateCustomerUseCase UpdateCustomerUseCase,
	deleteCustomerUseCase DeleteCustomerUseCase,
) *CustomerHandler {

	return &CustomerHandler{
		createCustomerUseCase: createCustomerUseCase,
		getCustomerUseCase:    getCustomerUseCase,
		listCustomersUseCase:  listCustomersUseCase,
		updateCustomerUseCase: updateCustomerUseCase,
		deleteCustomerUseCase: deleteCustomerUseCase,
	}
}

// Create
//
// @Summary Criar Cliente
// @Description Cadastra um novo cliente.
// @Tags Customers
// @Accept json
// @Produce json
// @Param request body customerdto.CreateCustomerRequest true "Cliente"
// @Success 201 {object} customer.CustomerResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Router /customers [post]
func (h *CustomerHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request customerdto.CreateCustomerRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	customer := mapper.ToCustomer(request)

	if err := h.createCustomerUseCase.Execute(customer); err != nil {

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
		mapper.ToCustomerResponse(customer),
	)
}

// GetByID
//
// @Summary Buscar Cliente
// @Description Busca um cliente pelo ID.
// @Tags Customers
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} customer.CustomerResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	customer, err := h.getCustomerUseCase.Execute(id)

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
		mapper.ToCustomerResponse(customer),
	)
}

// List
//
// @Summary Listar Clientes
// @Description Retorna todos os clientes.
// @Tags Customers
// @Produce json
// @Success 200 {array} customer.CustomerResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /customers [get]
func (h *CustomerHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	customers, err := h.listCustomersUseCase.Execute()

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
		mapper.ToCustomerResponseList(customers),
	)
}

// Update
//
// @Summary Atualizar Cliente
// @Description Atualiza um cliente.
// @Tags Customers
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body customer.UpdateCustomerRequest true "Cliente"
// @Success 200 {object} customer.CustomerResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /customers/{id} [put]
func (h *CustomerHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	var request customerdto.UpdateCustomerRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	customer := mapper.ToCustomerUpdate(
		id,
		request,
	)

	if err := h.updateCustomerUseCase.Execute(customer); err != nil {

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
		mapper.ToCustomerResponse(customer),
	)
}

// Delete
//
// @Summary Remover Cliente
// @Description Remove um cliente.
// @Tags Customers
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	if err := h.deleteCustomerUseCase.Execute(id); err != nil {

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
