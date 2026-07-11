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

	"desafio-go/infrastructure/http/httpx"

	customerdto "desafio-go/internal/dto/customer"
	"desafio-go/internal/mapper"

	customerusecase "desafio-go/internal/usecase/customer"
)

type CustomerHandler struct {
	createCustomerUseCase *customerusecase.CreateCustomerUseCase
	getCustomerUseCase    *customerusecase.GetCustomerUseCase
	listCustomersUseCase  *customerusecase.ListCustomersUseCase
	updateCustomerUseCase *customerusecase.UpdateCustomerUseCase
	deleteCustomerUseCase *customerusecase.DeleteCustomerUseCase
}

func NewCustomerHandler(
	createCustomerUseCase *customerusecase.CreateCustomerUseCase,
	getCustomerUseCase *customerusecase.GetCustomerUseCase,
	listCustomersUseCase *customerusecase.ListCustomersUseCase,
	updateCustomerUseCase *customerusecase.UpdateCustomerUseCase,
	deleteCustomerUseCase *customerusecase.DeleteCustomerUseCase,
) *CustomerHandler {

	return &CustomerHandler{
		createCustomerUseCase: createCustomerUseCase,
		getCustomerUseCase:    getCustomerUseCase,
		listCustomersUseCase:  listCustomersUseCase,
		updateCustomerUseCase: updateCustomerUseCase,
		deleteCustomerUseCase: deleteCustomerUseCase,
	}
}

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
