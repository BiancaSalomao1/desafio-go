package integration

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

Regras verificadas neste pacote:

✓ cliente obrigatório
✓ pedido deve possuir itens
✓ quantidade válida
✓ produtos duplicados não são permitidos
✓ transições de estado do pedido
✓ autenticação JWT
✓ persistência no PostgreSQL

Observação:
A validação de existência/estoque do produto pertence ao Product Service
e é exercida pela Saga por meio de eventos RabbitMQ. Esses cenários não
devem ser testados como respostas síncronas da Orders API.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	customerdto "orders-api/internal/dto/customer"
	orderdto "orders-api/internal/dto/order"
	productdto "orders-api/internal/dto/product"

	"github.com/stretchr/testify/assert"
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

func waitForOrderStatus(
	t *testing.T,
	ts *TestServer,
	token string,
	orderID string,
	expectedStatus string,
) {
	t.Helper()

	assert.Eventually(
		t,
		func() bool {
			req, err := authenticatedRequest(
				http.MethodGet,
				ts.Server.URL+"/orders/"+orderID,
				token,
				nil,
			)
			if err != nil {
				return false
			}

			resp := doRequest(t, req)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return false
			}

			order := decodeResponse[orderdto.OrderResponse](
				t,
				resp.Body,
			)

			return order.Status == expectedStatus
		},
		5*time.Second,
		100*time.Millisecond,
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
