/*
order_test

Responsabilidades:

- testar toda a API de Pedidos;
- validar regras de negócio;
- validar integração entre HTTP e PostgreSQL;
- validar movimentação de estoque.

Fluxo:

HTTP

	↓

JWT Middleware

	↓

Order Handler

	↓

Order UseCase

	↓

Repositories

	↓

# PostgreSQL

Regras verificadas:

✓ cliente obrigatório
✓ produto obrigatório
✓ quantidade válida
✓ estoque suficiente
✓ redução de estoque
*/
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	customerdto "desafio-go/orders-api/internal/dto/customer"
	orderdto "desafio-go/orders-api/internal/dto/order"
	productdto "desafio-go/orders-api/internal/dto/product"
)

func createOrder(
	t *testing.T,
	ts *TestServer,
	token string,
	request orderdto.CreateOrderRequest,
) orderdto.OrderResponse {

	t.Helper()

	body := mustMarshal(t, request)

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertCreated(t, resp)

	return decodeResponse[orderdto.OrderResponse](
		t,
		resp.Body,
	)
}

func TestCreateOrder(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "João",
			Email: "joao@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  2,
				},
			},
		},
	)

	if order.ID == "" {
		t.Fatal("expected id")
	}

	if order.CustomerID != customer.ID {
		t.Fatal("invalid customer")
	}

	if len(order.Items) != 1 {
		t.Fatal("expected one item")
	}

	if order.Status != "PENDING" {
		t.Fatal("expected pending status")
	}
}

func TestCreateOrderWithoutItems(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateOrderCustomerNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 80,
			Stock: 20,
		},
	)

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: "11111111-1111-1111-1111-111111111111",
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  1,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestCreateOrderProductNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		},
	)

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: "11111111-1111-1111-1111-111111111111",
				Quantity:  1,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestCreateOrderQuantityZero(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		})

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  0,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateOrderNegativeQuantity(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 120,
			Stock: 20,
		})

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  -1,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateOrderInsufficientStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Monitor",
			Price: 1500,
			Stock: 2,
		})

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  10,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertConflict(t, resp)
}

func TestCreateOrderMultipleProducts(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		})

	notebook := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 5,
		})

	mouse := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 80,
			Stock: 50,
		})

	order := createOrder(t, ts, token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: notebook.ID,
					Quantity:  1,
				},
				{
					ProductID: mouse.ID,
					Quantity:  3,
				},
			},
		})

	if len(order.Items) != 2 {
		t.Fatalf("expected 2 items got %d", len(order.Items))
	}
}

func TestCreateOrderReducesStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "José",
			Email: "jose@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "SSD",
			Price: 600,
			Stock: 10,
		})

	createOrder(t, ts, token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  4,
				},
			},
		})

	req, _ := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	updated := decodeResponse[productdto.ProductResponse](
		t,
		resp.Body,
	)

	if updated.Stock != 6 {
		t.Fatalf(
			"expected stock 6 got %d",
			updated.Stock,
		)
	}
}
func TestGetOrderByID(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "João",
			Email: "joao@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  2,
				},
			},
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders/"+order.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[orderdto.OrderResponse](
		t,
		resp.Body,
	)

	if response.ID != order.ID {
		t.Fatal("invalid id")
	}

	if response.CustomerID != customer.ID {
		t.Fatal("invalid customer")
	}

	if len(response.Items) != 1 {
		t.Fatal("expected one item")
	}
}

func TestGetOrderNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders/11111111-1111-1111-1111-111111111111",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestListOrdersEmpty(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	orders := decodeResponse[[]orderdto.OrderResponse](
		t,
		resp.Body,
	)

	if len(orders) != 0 {
		t.Fatalf("expected empty list got %d", len(orders))
	}
}

func TestListOrders(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 100,
			Stock: 100,
		},
	)

	for i := 0; i < 3; i++ {
		createOrder(
			t,
			ts,
			token,
			orderdto.CreateOrderRequest{
				CustomerID: customer.ID,
				Items: []orderdto.CreateOrderItemRequest{
					{
						ProductID: product.ID,
						Quantity:  1,
					},
				},
			},
		)
	}

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	orders := decodeResponse[[]orderdto.OrderResponse](
		t,
		resp.Body,
	)

	if len(orders) != 3 {
		t.Fatalf("expected 3 orders got %d", len(orders))
	}
}

func TestListOrdersPagination(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "SSD",
			Price: 600,
			Stock: 100,
		},
	)

	for i := 0; i < 10; i++ {
		createOrder(
			t,
			ts,
			token,
			orderdto.CreateOrderRequest{
				CustomerID: customer.ID,
				Items: []orderdto.CreateOrderItemRequest{
					{
						ProductID: product.ID,
						Quantity:  1,
					},
				},
			},
		)
	}

	req, err := authenticatedRequest(
		http.MethodGet,
		fmt.Sprintf("%s/orders?page=1&limit=5", ts.Server.URL),
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	orders := decodeResponse[[]orderdto.OrderResponse](
		t,
		resp.Body,
	)

	if len(orders) != 5 {
		t.Fatalf("expected 5 orders got %d", len(orders))
	}
}
func TestPayOrder(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  2,
				},
			},
		},
	)

	req, err := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNoContent(t, resp)

}

func TestPayOrderTwice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 80,
			Stock: 20,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  1,
				},
			},
		},
	)

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertConflict(t, resp)
}

func TestPayOrderNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/11111111-1111-1111-1111-111111111111/pay",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestPayOrderStatusTransition(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Teclado",
			Price: 300,
			Stock: 5,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  1,
				},
			},
		},
	)

	if order.Status != "PENDING" {
		t.Fatalf(
			"expected PENDING got %s",
			order.Status,
		)
	}

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertConflict(t, resp)

	req, _ = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders/"+order.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	updated := decodeResponse[orderdto.OrderResponse](
		t,
		resp.Body,
	)

	if updated.Status != "PAID" {
		t.Fatalf(
			"expected PAID got %s",
			updated.Status,
		)
	}
}
func TestCancelOrder(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  2,
				},
			},
		},
	)

	req, err := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNoContent(t, resp)

}

func TestCancelOrderTwice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 120,
			Stock: 20,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  1,
				},
			},
		},
	)

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertConflict(t, resp)
}

func TestCancelPaidOrder(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Teclado",
			Price: 300,
			Stock: 5,
		},
	)

	order := createOrder(
		t,
		ts,
		token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  1,
				},
			},
		},
	)

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertConflict(t, resp)
}

func TestCancelOrderNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/11111111-1111-1111-1111-111111111111/cancel",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestCancelOrderRestoresStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		})

	order := createOrder(t, ts, token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  4,
				},
			},
		})

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	updated := decodeResponse[productdto.ProductResponse](t, resp.Body)

	if updated.Stock != 10 {
		t.Fatalf("expected stock 10 got %d", updated.Stock)
	}
}

func TestPaidOrderDoesNotRestoreStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Monitor",
			Price: 2000,
			Stock: 8,
		})

	order := createOrder(t, ts, token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  3,
				},
			},
		})

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/pay",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	productResponse := decodeResponse[productdto.ProductResponse](t, resp.Body)

	if productResponse.Stock != 5 {
		t.Fatalf("expected stock 5 got %d", productResponse.Stock)
	}
}

func TestCancelOrderRestoresMultipleProducts(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	productA := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		})

	productB := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 100,
			Stock: 20,
		})

	order := createOrder(t, ts, token,
		orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: productA.ID,
					Quantity:  2,
				},
				{
					ProductID: productB.ID,
					Quantity:  5,
				},
			},
		})

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/"+order.ID+"/cancel",
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	for _, tc := range []struct {
		id    string
		stock int
	}{
		{productA.ID, 10},
		{productB.ID, 20},
	} {

		req, _ = authenticatedRequest(
			http.MethodGet,
			ts.Server.URL+"/products/"+tc.id,
			token,
			nil,
		)

		resp = doRequest(t, req)

		product := decodeResponse[productdto.ProductResponse](t, resp.Body)
		resp.Body.Close()

		if product.Stock != tc.stock {
			t.Fatalf(
				"expected stock %d got %d",
				tc.stock,
				product.Stock,
			)
		}
	}
}

func TestStockNeverBecomesNegative(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(t, ts, token,
		customerdto.CreateCustomerRequest{
			Name:  "Cliente",
			Email: "cliente@email.com",
		})

	product := createProduct(t, ts, token,
		productdto.CreateProductRequest{
			Name:  "SSD",
			Price: 500,
			Stock: 2,
		})

	body := mustMarshal(t, orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  3,
			},
		},
	})

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertConflict(t, resp)

	req, _ = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	current := decodeResponse[productdto.ProductResponse](t, resp.Body)

	if current.Stock != 2 {
		t.Fatalf(
			"stock should remain 2 got %d",
			current.Stock,
		)
	}
}
func TestCreateOrderWithoutToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	body := mustMarshal(t, orderdto.CreateOrderRequest{})

	req, err := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestCreateOrderInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	body := mustMarshal(t, orderdto.CreateOrderRequest{})

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		"invalid-token",
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestCreateOrderInvalidJSON(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBufferString(`{"customerId":`),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestPayOrderWithoutToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/123/pay",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestCancelOrderWithoutToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/123/cancel",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}
func TestListOrdersUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/orders",
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestGetOrderUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/orders/123",
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestPayOrderInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/123/pay",
		"invalid-token",
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestCancelOrderInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := authenticatedRequest(
		http.MethodPatch,
		ts.Server.URL+"/orders/123/cancel",
		"invalid-token",
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestGetOrderInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/orders/123",
		"invalid-token",
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestCreateOrderDuplicatedProduct(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:     "Cliente",
			Email:    "cliente@email.com",
			Password: "123456",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 10,
		},
	)

	request := orderdto.CreateOrderRequest{
		CustomerID: customer.ID,
		Items: []orderdto.CreateOrderItemRequest{
			{
				ProductID: product.ID,
				Quantity:  2,
			},
			{
				ProductID: product.ID,
				Quantity:  3,
			},
		},
	}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/orders",
		token,
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf(
			"expected %d got %d",
			http.StatusConflict,
			resp.StatusCode,
		)
	}
}

func TestCreateOrderConcurrentStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:     "Cliente Teste",
			Email:    "cliente@teste.com",
			Password: "123456",
		},
	)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 5,
		},
	)
	var wg sync.WaitGroup

	statuses := make(chan int, 2)

	create := func() {
		defer wg.Done()

		request := orderdto.CreateOrderRequest{
			CustomerID: customer.ID,
			Items: []orderdto.CreateOrderItemRequest{
				{
					ProductID: product.ID,
					Quantity:  5,
				},
			},
		}

		body, _ := json.Marshal(request)

		req, _ := authenticatedRequest(
			http.MethodPost,
			ts.Server.URL+"/orders",
			token,
			bytes.NewBuffer(body),
		)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()

		statuses <- resp.StatusCode
	}

	wg.Add(2)

	go create()
	go create()

	wg.Wait()
	close(statuses)

	var created int
	var failed int

	for status := range statuses {
		switch status {
		case http.StatusCreated:
			created++
		default:
			failed++
		}
	}

	if created != 1 {
		t.Fatalf("expected 1 created order, got %d", created)
	}

	if failed != 1 {
		t.Fatalf("expected 1 failed order, got %d", failed)
	}
}
