package handler

/*
struct OrderHandler

Responsabilidades:
- receber requisições HTTP relacionadas a pedidos;
- converter DTOs para entidades de domínio;
- executar os casos de uso;
- retornar respostas HTTP em JSON.

Campos:
- createOrderUseCase
- getOrderUseCase
- payOrderUseCase
- cancelOrderUseCase

Métodos:
- NewOrderHandler()
- Create()
- GetByID()
- Pay()
- Cancel()
*/

import (
	"net/http"

	"desafio-go/infrastructure/http/httpx"

	orderdto "desafio-go/internal/dto/order"
	"desafio-go/internal/mapper"

	orderusecase "desafio-go/internal/usecase/order"
)

type OrderHandler struct {
	createOrderUseCase *orderusecase.CreateOrderUseCase
	getOrderUseCase    *orderusecase.GetOrderUseCase
	payOrderUseCase    *orderusecase.PayOrderUseCase
	cancelOrderUseCase *orderusecase.CancelOrderUseCase
}

func NewOrderHandler(
	createOrderUseCase *orderusecase.CreateOrderUseCase,
	getOrderUseCase *orderusecase.GetOrderUseCase,
	payOrderUseCase *orderusecase.PayOrderUseCase,
	cancelOrderUseCase *orderusecase.CancelOrderUseCase,
) *OrderHandler {

	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
		getOrderUseCase:    getOrderUseCase,
		payOrderUseCase:    payOrderUseCase,
		cancelOrderUseCase: cancelOrderUseCase,
	}
}

func (h *OrderHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request orderdto.CreateOrderRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	order := mapper.ToOrder(request)

	if err := h.createOrderUseCase.Execute(order); err != nil {

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
		mapper.ToOrderResponse(order),
	)
}

func (h *OrderHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	order, err := h.getOrderUseCase.Execute(id)

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
		mapper.ToOrderResponse(order),
	)
}

func (h *OrderHandler) Pay(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	if err := h.payOrderUseCase.Execute(id); err != nil {

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

func (h *OrderHandler) Cancel(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	if err := h.cancelOrderUseCase.Execute(id); err != nil {

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
