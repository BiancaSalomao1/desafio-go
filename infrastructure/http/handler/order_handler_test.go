package handler

import (
	"bytes"
	"desafio-go/infrastructure/http/testutil"
	"desafio-go/internal/domain"
	orderdto "desafio-go/internal/dto/order"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =========================
// Mocks
// =========================

type createOrderUseCaseMock struct {
	execute func(*domain.Order) error
}

func (m *createOrderUseCaseMock) Execute(order *domain.Order) error {
	if m.execute != nil {
		return m.execute(order)
	}
	return nil
}

type getOrderUseCaseMock struct {
	execute func(string) (*domain.Order, error)
}

func (m *getOrderUseCaseMock) Execute(id string) (*domain.Order, error) {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil, nil
}

type listOrdersUseCaseMock struct {
	execute func(int, int) ([]*domain.Order, error)
}

func (m *listOrdersUseCaseMock) Execute(limit int, offset int) ([]*domain.Order, error) {
	if m.execute != nil {
		return m.execute(limit, offset)
	}
	return []*domain.Order{}, nil
}

type payOrderUseCaseMock struct {
	execute func(string) error
}

func (m *payOrderUseCaseMock) Execute(id string) error {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil
}

type cancelOrderUseCaseMock struct {
	execute func(string) error
}

func (m *cancelOrderUseCaseMock) Execute(id string) error {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil
}

// =========================
// Helpers
// =========================

func newOrderHandler(
	create CreateOrderUseCase,
	get GetOrderUseCase,
	list ListOrdersUseCase,
	pay PayOrderUseCase,
	cancel CancelOrderUseCase,
) *OrderHandler {
	return NewOrderHandler(
		create,
		get,
		list,
		pay,
		cancel,
	)
}

func validOrder() *domain.Order {
	return &domain.Order{
		ID:         "ORD001",
		CustomerID: "UI001",
		Items: []domain.OrderItem{
			{
				ProductID: "PROD001",
				Quantity:  1,
			},
		},

		Status: domain.OrderStatusPending,
	}
}

func TestOrderHandler_Create(t *testing.T) {

	t.Run("should create order", func(t *testing.T) {

		mock := &createOrderUseCaseMock{
			execute: func(order *domain.Order) error {

				if order.CustomerID != "UI001" {
					t.Fatalf("unexpected customer_id: %s", order.CustomerID)
				}

				if order.Items[0].Quantity != 1 {
					t.Fatalf("unexpected quantity: %d", order.Items[0].Quantity)
				}

				return nil
			},
		}

		handler := newOrderHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := orderdto.CreateOrderRequest{
			CustomerID: "UI001",
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: "PROD001",
					Quantity:  1,
				},
			},
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/orders",
			requestBody,
		)

		handler.Create(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusCreated,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[orderdto.OrderResponse](
			t,
			rec,
		)

		if response.CustomerID != "UI001" {
			t.Fatalf("expected UI001, got %s", response.CustomerID)
		}

		if response.Items[0].Quantity != 1 {
			t.Fatalf("expected 1, got %d", response.Items[0].Quantity)
		}
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newOrderHandler(
			&createOrderUseCaseMock{},
			nil,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/orders",
			bytes.NewBufferString("{invalid"),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &createOrderUseCaseMock{
			execute: func(order *domain.Order) error {
				return errors.New("customer not found")
			},
		}

		handler := newOrderHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := orderdto.CreateOrderRequest{
			CustomerID: "UI001",
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: "PROD001",
					Quantity:  1,
				},
			},
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/orders",
			requestBody,
		)

		handler.Create(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"customer not found",
		)
	})
}

func TestOrderHandler_GetByID(t *testing.T) {

	t.Run("should return order", func(t *testing.T) {

		mock := &getOrderUseCaseMock{
			execute: func(id string) (*domain.Order, error) {

				if id != "ORD001" {
					t.Fatalf("expected ORD001, got %s", id)
				}

				return validOrder(), nil
			},
		}

		handler := newOrderHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/ORD001",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[orderdto.OrderResponse](
			t,
			rec,
		)

		if response.ID != "ORD001" {
			t.Fatalf("expected ORD001, got %s", response.ID)
		}

		if response.CustomerID != "UI001" {
			t.Fatalf("expected UI001, got %s", response.CustomerID)
		}

		if response.Items[0].Quantity != 1 {
			t.Fatalf("expected 1, got %d", response.Items[0].Quantity)
		}
	})

	t.Run("should return not found", func(t *testing.T) {

		mock := &getOrderUseCaseMock{
			execute: func(id string) (*domain.Order, error) {
				return nil, errors.New("order not found")
			},
		}

		handler := newOrderHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders/ORD001",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNotFound,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"order not found",
		)
	})
}

func TestOrderHandler_Pay(t *testing.T) {

	t.Run("should pay order", func(t *testing.T) {

		mock := &payOrderUseCaseMock{
			execute: func(id string) error {

				if id != "ORD001" {
					t.Fatalf("expected ORD001, got %s", id)
				}

				return nil
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/orders/ORD001/pay",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.Pay(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNoContent,
		)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &payOrderUseCaseMock{
			execute: func(id string) error {
				return errors.New("order not found")
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/orders/ORD001/pay",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.Pay(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"order not found",
		)
	})
}

func TestOrderHandler_Cancel(t *testing.T) {

	t.Run("should cancel order", func(t *testing.T) {

		mock := &cancelOrderUseCaseMock{
			execute: func(id string) error {

				if id != "ORD001" {
					t.Fatalf("expected ORD001, got %s", id)
				}

				return nil
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/orders/ORD001/cancel",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.Cancel(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNoContent,
		)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &cancelOrderUseCaseMock{
			execute: func(id string) error {
				return errors.New("order not found")
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/orders/ORD001/cancel",
			nil,
		)

		req.SetPathValue("id", "ORD001")

		rec := httptest.NewRecorder()

		handler.Cancel(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"order not found",
		)
	})
}

func TestOrderHandler_List(t *testing.T) {

	t.Run("should list orders", func(t *testing.T) {

		mock := &listOrdersUseCaseMock{
			execute: func(limit int, offset int) ([]*domain.Order, error) {

				return []*domain.Order{
					validOrder(),
				}, nil
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders?limit=10&offset=0",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[[]orderdto.OrderResponse](
			t,
			rec,
		)

		if len(response) != 1 {
			t.Fatalf("expected 1 order, got %d", len(response))
		}

		if response[0].CustomerID != "UI001" {
			t.Fatalf("expected UI001, got %s", response[0].CustomerID)
		}
	})

	t.Run("should return internal server error", func(t *testing.T) {

		mock := &listOrdersUseCaseMock{
			execute: func(limit int, offset int) ([]*domain.Order, error) {
				return nil, errors.New("database error")
			},
		}

		handler := newOrderHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/orders?limit=10&offset=0",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusInternalServerError,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"database error",
		)

		t.Run("should use default limit when invalid", func(t *testing.T) {

			mock := &listOrdersUseCaseMock{
				execute: func(limit, offset int) ([]*domain.Order, error) {

					if limit != 10 {
						t.Fatalf("expected limit 10, got %d", limit)
					}

					if offset != 0 {
						t.Fatalf("expected offset 0, got %d", offset)
					}

					return []*domain.Order{}, nil
				},
			}

			handler := newOrderHandler(
				nil,
				nil,
				mock,
				nil,
				nil,
			)

			req := httptest.NewRequest(
				http.MethodGet,
				"/orders?limit=abc",
				nil,
			)

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			testutil.AssertStatus(t, rec, http.StatusOK)
		})
		t.Run("should cap limit at 100", func(t *testing.T) {

			mock := &listOrdersUseCaseMock{
				execute: func(limit, offset int) ([]*domain.Order, error) {

					if limit != 100 {
						t.Fatalf("expected limit 100, got %d", limit)
					}

					return []*domain.Order{}, nil
				},
			}

			handler := newOrderHandler(nil, nil, mock, nil, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/orders?limit=500",
				nil,
			)

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			testutil.AssertStatus(t, rec, http.StatusOK)
		})
		t.Run("should use offset zero when invalid", func(t *testing.T) {

			mock := &listOrdersUseCaseMock{
				execute: func(limit, offset int) ([]*domain.Order, error) {

					if offset != 0 {
						t.Fatalf("expected offset 0, got %d", offset)
					}

					return []*domain.Order{}, nil
				},
			}

			handler := newOrderHandler(nil, nil, mock, nil, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/orders?offset=abc",
				nil,
			)

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			testutil.AssertStatus(t, rec, http.StatusOK)
		})
		t.Run("should use offset zero when negative", func(t *testing.T) {

			mock := &listOrdersUseCaseMock{
				execute: func(limit, offset int) ([]*domain.Order, error) {

					if offset != 0 {
						t.Fatalf("expected offset 0, got %d", offset)
					}

					return []*domain.Order{}, nil
				},
			}

			handler := newOrderHandler(nil, nil, mock, nil, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/orders?offset=-5",
				nil,
			)

			rec := httptest.NewRecorder()

			handler.List(rec, req)

			testutil.AssertStatus(t, rec, http.StatusOK)
		})

	})
}
