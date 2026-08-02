package handler

import (
	"bytes"
	"desafio-go/infrastructure/http/testutil"
	"desafio-go/internal/domain"
	customerdto "desafio-go/internal/dto/customer"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =========================
// Mocks
// =========================

type createCustomerUseCaseMock struct {
	execute func(*domain.Customer) error
}

func (m *createCustomerUseCaseMock) Execute(customer *domain.Customer) error {
	if m.execute != nil {
		return m.execute(customer)
	}
	return nil
}

type getCustomerUseCaseMock struct {
	execute func(string) (*domain.Customer, error)
}

func (m *getCustomerUseCaseMock) Execute(id string) (*domain.Customer, error) {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil, nil
}

type listCustomersUseCaseMock struct {
	execute func() ([]*domain.Customer, error)
}

func (m *listCustomersUseCaseMock) Execute() ([]*domain.Customer, error) {
	if m.execute != nil {
		return m.execute()
	}
	return []*domain.Customer{}, nil
}

type updateCustomerUseCaseMock struct {
	execute func(*domain.Customer) error
}

func (m *updateCustomerUseCaseMock) Execute(customer *domain.Customer) error {
	if m.execute != nil {
		return m.execute(customer)
	}
	return nil
}

type deleteCustomerUseCaseMock struct {
	execute func(string) error
}

func (m *deleteCustomerUseCaseMock) Execute(id string) error {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil
}

// =========================
// Helpers
// =========================

func newCustomerHandler(
	create CreateCustomerUseCase,
	get GetCustomerUseCase,
	list ListCustomersUseCase,
	update UpdateCustomerUseCase,
	delete DeleteCustomerUseCase,
) *CustomerHandler {

	return NewCustomerHandler(
		create,
		get,
		list,
		update,
		delete,
	)
}

func validCustomer() *domain.Customer {
	return domain.NewCustomer(
		"UI001",
		"John Doe",
		"john.doe@example.com",
		"hash",
	)
}

func TestCustomerHandler_Create(t *testing.T) {

	t.Run("should create customer", func(t *testing.T) {

		mock := &createCustomerUseCaseMock{
			execute: func(customer *domain.Customer) error {

				if customer.Name != "John Doe" {
					t.Fatalf("unexpected name: %s", customer.Name)
				}

				if customer.Email != "john.doe@example.com" {
					t.Fatalf("unexpected email: %s", customer.Email)
				}

				return nil
			},
		}

		handler := newCustomerHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := customerdto.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/customers",
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

		response := testutil.DecodeResponse[customerdto.CustomerResponse](
			t,
			rec,
		)

		if response.Name != "John Doe" {
			t.Fatalf("expected John Doe, got %s", response.Name)
		}

		if response.Email != "john.doe@example.com" {
			t.Fatalf("expected john.doe@example.com, got %s", response.Email)
		}
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newCustomerHandler(
			&createCustomerUseCaseMock{},
			nil,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/customers",
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

		mock := &createCustomerUseCaseMock{
			execute: func(customer *domain.Customer) error {
				return errors.New("customer already exists")
			},
		}

		handler := newCustomerHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := customerdto.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john.doe@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/customers",
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
			"customer already exists",
		)
	})
}

func TestCustomerHandler_GetByID(t *testing.T) {

	t.Run("should return customer", func(t *testing.T) {

		mock := &getCustomerUseCaseMock{
			execute: func(id string) (*domain.Customer, error) {

				if id != "UI001" {
					t.Fatalf("expected UI001, got %s", id)
				}

				return validCustomer(), nil
			},
		}

		handler := newCustomerHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/customers/UI001",
			nil,
		)

		req.SetPathValue("id", "UI001")

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

		response := testutil.DecodeResponse[customerdto.CustomerResponse](
			t,
			rec,
		)

		if response.ID != "UI001" {
			t.Fatalf("expected UI001, got %s", response.ID)
		}

		if response.Name != "John Doe" {
			t.Fatalf("expected John Doe, got %s", response.Name)
		}

		if response.Email != "john.doe@example.com" {
			t.Fatalf("expected john.doe@example.com, got %s", response.Email)
		}
	})

	t.Run("should return not found", func(t *testing.T) {

		mock := &getCustomerUseCaseMock{
			execute: func(id string) (*domain.Customer, error) {
				return nil, errors.New("customer not found")
			},
		}

		handler := newCustomerHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/customers/UI001",
			nil,
		)

		req.SetPathValue("id", "UI001")

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
			"customer not found",
		)
	})
}

func TestCustomerHandler_List(t *testing.T) {

	t.Run("should list customers", func(t *testing.T) {

		mock := &listCustomersUseCaseMock{
			execute: func() ([]*domain.Customer, error) {

				return []*domain.Customer{
					validCustomer(),
				}, nil
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/customers",
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

		response := testutil.DecodeResponse[[]customerdto.CustomerResponse](
			t,
			rec,
		)

		if len(response) != 1 {
			t.Fatalf("expected 1 customer, got %d", len(response))
		}

		if response[0].Name != "John Doe" {
			t.Fatalf("expected John Doe, got %s", response[0].Name)
		}
	})

	t.Run("should return internal server error", func(t *testing.T) {

		mock := &listCustomersUseCaseMock{
			execute: func() ([]*domain.Customer, error) {
				return nil, errors.New("database error")
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/customers",
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
	})
}

func TestCustomerHandler_Update(t *testing.T) {

	t.Run("should update customer", func(t *testing.T) {

		mock := &updateCustomerUseCaseMock{
			execute: func(customer *domain.Customer) error {

				if customer.ID != "UI001" {
					t.Fatalf("expected UI001, got %s", customer.ID)
				}

				if customer.Name != "John Updated" {
					t.Fatalf("expected John Updated, got %s", customer.Name)
				}

				if customer.Email != "john.updated@example.com" {
					t.Fatalf("expected john.updated@example.com, got %s", customer.Email)
				}

				return nil
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		requestBody := customerdto.UpdateCustomerRequest{
			Name:  "John Updated",
			Email: "john.updated@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/customers/UI001",
			requestBody,
		)

		req.SetPathValue("id", "UI001")

		handler.Update(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[customerdto.CustomerResponse](
			t,
			rec,
		)

		if response.ID != "UI001" {
			t.Fatalf("expected UI001, got %s", response.ID)
		}

		if response.Name != "John Updated" {
			t.Fatalf("expected John Updated, got %s", response.Name)
		}

		if response.Email != "john.updated@example.com" {
			t.Fatalf("expected john.updated@example.com, got %s", response.Email)
		}
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newCustomerHandler(
			nil,
			nil,
			nil,
			&updateCustomerUseCaseMock{},
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPut,
			"/customers/UI001",
			bytes.NewBufferString("{invalid"),
		)

		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", "UI001")

		rec := httptest.NewRecorder()

		handler.Update(rec, req)

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

		mock := &updateCustomerUseCaseMock{
			execute: func(customer *domain.Customer) error {
				return errors.New("customer not found")
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		requestBody := customerdto.UpdateCustomerRequest{
			Name:  "John Updated",
			Email: "john.updated@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/customers/UI001",
			requestBody,
		)

		req.SetPathValue("id", "UI001")

		handler.Update(rec, req)

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

func TestCustomerHandler_Delete(t *testing.T) {

	t.Run("should delete customer", func(t *testing.T) {

		mock := &deleteCustomerUseCaseMock{
			execute: func(id string) error {

				if id != "UI001" {
					t.Fatalf("expected UI001, got %s", id)
				}

				return nil
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/customers/UI001",
			nil,
		)

		req.SetPathValue("id", "UI001")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNoContent,
		)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &deleteCustomerUseCaseMock{
			execute: func(id string) error {
				return errors.New("customer not found")
			},
		}

		handler := newCustomerHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/customers/UI001",
			nil,
		)

		req.SetPathValue("id", "UI001")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

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
