package handler

/*
- mocks
- newHandler()
- TestProductHandler_Create()
    - success
    - invalid json
    - empty body
    - usecase error

- TestProductHandler_List()
    - success
    - empty
    - error

- TestProductHandler_GetByID()
    - success
    - not found

- TestProductHandler_Update()
    - success
    - invalid json
    - usecase error

- TestProductHandler_Delete()
    - success
    - error

*/

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"desafio-go/infrastructure/http/testutil"
	"desafio-go/internal/domain"
	productdto "desafio-go/internal/dto/product"
)

// mocks
type createProductUseCaseMock struct {
	execute func(*domain.Product) error
}

func (m *createProductUseCaseMock) Execute(product *domain.Product) error {
	if m.execute == nil {
		return nil
	}

	return m.execute(product)
}

type getProductUseCaseMock struct {
	execute func(string) (*domain.Product, error)
}

func (m *getProductUseCaseMock) Execute(id string) (*domain.Product, error) {
	if m.execute == nil {
		return nil, nil
	}

	return m.execute(id)
}

type listProductsUseCaseMock struct {
	execute func() ([]*domain.Product, error)
}

func (m *listProductsUseCaseMock) Execute() ([]*domain.Product, error) {

	if m.execute == nil {
		return []*domain.Product{}, nil
	}

	return m.execute()
}

type updateProductUseCaseMock struct {
	execute func(*domain.Product) error
}

func (m *updateProductUseCaseMock) Execute(product *domain.Product) error {

	if m.execute == nil {
		return nil
	}

	return m.execute(product)
}

type deleteProductUseCaseMock struct {
	execute func(string) error
}

func (m *deleteProductUseCaseMock) Execute(id string) error {

	if m.execute == nil {
		return nil
	}

	return m.execute(id)
}

// helpers
func newHandler(
	create CreateProductUseCase,
	get GetProductUseCase,
	list ListProductsUseCase,
	update UpdateProductUseCase,
	delete DeleteProductUseCase,
) *ProductHandler {

	return NewProductHandler(
		create,
		get,
		list,
		update,
		delete,
	)
}

func validProduct() *domain.Product {

	return domain.NewProduct(
		"1",
		"Notebook",
		3500,
		10,
	)
}

// tests
func TestProductHandler_Create(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mock := &createProductUseCaseMock{
			execute: func(product *domain.Product) error {

				if product.Name != "Notebook" {
					t.Fatal("invalid name")
				}

				if product.Price != 3500 {
					t.Fatal("invalid price")
				}

				if product.Stock != 10 {
					t.Fatal("invalid stock")
				}

				return nil
			},
		}

		handler := newHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/products",
			productdto.CreateProductRequest{
				Name:  "Notebook",
				Price: 3500,
				Stock: 10,
			},
		)

		handler.Create(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusCreated,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[productdto.ProductResponse](
			t,
			rec,
		)

		if response.Name != "Notebook" {
			t.Fatal("invalid response")
		}
	})

	t.Run("invalid json", func(t *testing.T) {

		handler := newHandler(
			&createProductUseCaseMock{},
			nil,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/products",
			bytes.NewBufferString("{"),
		)

		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)
	})
	t.Run("empty body", func(t *testing.T) {

		handler := newHandler(
			&createProductUseCaseMock{},
			nil,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/products",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.Create(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)
	})

	t.Run("use case error", func(t *testing.T) {

		mock := &createProductUseCaseMock{
			execute: func(product *domain.Product) error {
				return errors.New("create error")
			},
		}

		handler := newHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/products",
			productdto.CreateProductRequest{
				Name:  "Notebook",
				Price: 3500,
				Stock: 10,
			},
		)

		handler.Create(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"create error",
		)
	})
}

func TestProductHandler_List(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mock := &listProductsUseCaseMock{
			execute: func() ([]*domain.Product, error) {

				return []*domain.Product{
					validProduct(),
				}, nil
			},
		}

		handler := newHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/products",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[[]productdto.ProductResponse](t, rec)

		if len(response) != 1 {
			t.Fatalf("expected 1 product, got %d", len(response))
		}

		if response[0].Name != "Notebook" {
			t.Fatal("invalid response")
		}
	})

	t.Run("empty list", func(t *testing.T) {

		mock := &listProductsUseCaseMock{
			execute: func() ([]*domain.Product, error) {
				return []*domain.Product{}, nil
			},
		}

		handler := newHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/products",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		response := testutil.DecodeResponse[[]productdto.ProductResponse](t, rec)

		if len(response) != 0 {
			t.Fatal("expected empty list")
		}
	})

	t.Run("use case error", func(t *testing.T) {

		mock := &listProductsUseCaseMock{
			execute: func() ([]*domain.Product, error) {
				return nil, errors.New("database error")
			},
		}

		handler := newHandler(
			nil,
			nil,
			mock,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/products",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusInternalServerError,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"database error",
		)
	})
}

func TestProductHandler_GetByID(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mock := &getProductUseCaseMock{
			execute: func(id string) (*domain.Product, error) {

				if id != "1" {
					t.Fatalf("expected id 1, got %s", id)
				}

				return validProduct(), nil
			},
		}

		handler := newHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/products/1",
			nil,
		)

		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		handler.GetByID(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[productdto.ProductResponse](t, rec)

		if response.ID != "1" {
			t.Fatal("invalid id")
		}

		if response.Name != "Notebook" {
			t.Fatal("invalid name")
		}
	})

	t.Run("not found", func(t *testing.T) {

		mock := &getProductUseCaseMock{
			execute: func(id string) (*domain.Product, error) {
				return nil, errors.New("product not found")
			},
		}

		handler := newHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/products/999",
			nil,
		)

		req.SetPathValue("id", "999")

		rec := httptest.NewRecorder()

		handler.GetByID(
			rec,
			req,
		)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNotFound,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"product not found",
		)
	})
}

func TestProductHandler_Update(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mock := &updateProductUseCaseMock{
			execute: func(product *domain.Product) error {

				if product.ID != "1" {
					t.Fatalf("expected id 1, got %s", product.ID)
				}

				if product.Name != "Notebook Gamer" {
					t.Fatalf("unexpected name")
				}

				if product.Price != 5000 {
					t.Fatalf("unexpected price")
				}

				if product.Stock != 20 {
					t.Fatalf("unexpected stock")
				}

				return nil
			},
		}

		handler := newHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/products/1",
			productdto.UpdateProductRequest{
				Name:  "Notebook Gamer",
				Price: 5000,
				Stock: 20,
			},
		)

		req.SetPathValue("id", "1")

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

		response := testutil.DecodeResponse[productdto.ProductResponse](
			t,
			rec,
		)

		if response.ID != "1" {
			t.Fatal("invalid id")
		}

		if response.Name != "Notebook Gamer" {
			t.Fatal("invalid response")
		}
	})

	t.Run("invalid json", func(t *testing.T) {

		handler := newHandler(
			nil,
			nil,
			nil,
			&updateProductUseCaseMock{},
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPut,
			"/products/1",
			bytes.NewBufferString("{"),
		)

		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)
	})

	t.Run("use case error", func(t *testing.T) {

		mock := &updateProductUseCaseMock{
			execute: func(product *domain.Product) error {
				return errors.New("update error")
			},
		}

		handler := newHandler(
			nil,
			nil,
			nil,
			mock,
			nil,
		)

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/products/1",
			productdto.UpdateProductRequest{
				Name:  "Notebook",
				Price: 100,
				Stock: 5,
			},
		)

		req.SetPathValue("id", "1")

		handler.Update(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"update error",
		)
	})
}
func TestProductHandler_Delete(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mock := &deleteProductUseCaseMock{
			execute: func(id string) error {

				if id != "1" {
					t.Fatalf("expected id 1, got %s", id)
				}

				return nil
			},
		}

		handler := newHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/products/1",
			nil,
		)

		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusNoContent,
		)
	})

	t.Run("use case error", func(t *testing.T) {

		mock := &deleteProductUseCaseMock{
			execute: func(id string) error {
				return errors.New("delete error")
			},
		}

		handler := newHandler(
			nil,
			nil,
			nil,
			nil,
			mock,
		)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/products/1",
			nil,
		)

		req.SetPathValue("id", "1")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusBadRequest,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"delete error",
		)
	})
}
