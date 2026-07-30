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
- listOrderUseCase

Métodos:
- NewOrderHandler()
- Create()
- GetByID()
- Pay()
- Cancel()
- List()
*/

import (
	"net/http"

	"desafio-go/infrastructure/http/httpx"

	orderdto "desafio-go/internal/dto/order"
	"desafio-go/internal/mapper"

	"strconv"
)

type OrderHandler struct {
	createOrderUseCase CreateOrderUseCase
	getOrderUseCase    GetOrderUseCase
	listOrdersUseCase  ListOrdersUseCase
	payOrderUseCase    PayOrderUseCase
	cancelOrderUseCase CancelOrderUseCase
}

func NewOrderHandler(
	create CreateOrderUseCase,
	get GetOrderUseCase,
	list ListOrdersUseCase,
	pay PayOrderUseCase,
	cancel CancelOrderUseCase,
) *OrderHandler {

	return &OrderHandler{
		createOrderUseCase: create,
		getOrderUseCase:    get,
		listOrdersUseCase:  list,
		payOrderUseCase:    pay,
		cancelOrderUseCase: cancel,
	}
}

// Create
//
// @Summary Criar Pedido
// @Description Cria um novo pedido.
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body order.CreateOrderRequest true "Pedido"
// @Success 201 {object} order.OrderResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /orders [post]
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

// GetByID
//
// @Summary Buscar Pedido
// @Description Busca um pedido pelo ID.
// @Tags Orders
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} order.OrderResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id} [get]
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

// Pay
//
// @Summary Pagar Pedido
// @Description Altera o status do pedido para PAID.
// @Tags Orders
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id}/pay [patch]
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

// Cancel
//
// @Summary Cancelar Pedido
// @Description Cancela um pedido e devolve o estoque.
// @Tags Orders
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Failure 400 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /orders/{id}/cancel [patch]
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

// List
//
// @Summary Listar Pedidos
// @Description Retorna todos os pedidos.
// @Tags Orders
// @Produce json
// @Success 200 {array} order.OrderResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Security BearerAuth
// @Router /orders [get]
// @Param limit query int false "Quantidade máxima de registros"
// @Param offset query int false "Quantidade de registros ignorados"
func (h *OrderHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	orders, err := h.listOrdersUseCase.Execute(limit, offset)

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
		mapper.ToOrderResponseList(orders),
	)
}
