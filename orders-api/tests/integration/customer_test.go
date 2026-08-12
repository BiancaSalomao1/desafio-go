/*
customer_test

Responsabilidades:

- testar toda a API de Clientes;
- validar integração entre HTTP e PostgreSQL;
- validar autenticação JWT;
- validar regras de negócio.

Fluxo testado:

HTTP
 ↓
JWT Middleware
 ↓
Router
 ↓
Handler
 ↓
UseCase
 ↓
Repository PostgreSQL
 ↓
Banco

Cenários desta parte:

✓ Create Customer
✓ Create Invalid JSON
✓ Create Invalid Name
✓ Create Invalid Email
✓ Create Without Token
*/

package integration

import (
	"bytes"
	"net/http"
	"testing"

	customerdto "desafio-go/orders-api/internal/dto/customer"
)

func createCustomer(
	t *testing.T,
	ts *TestServer,
	token string,
	request customerdto.CreateCustomerRequest,
) customerdto.CustomerResponse {

	t.Helper()

	if request.Password == "" {
		request.Password = "123456"
	}

	body := mustMarshal(t, request)

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertCreated(t, resp)

	return decodeResponse[customerdto.CustomerResponse](
		t,
		resp.Body,
	)
}

func TestCreateCustomer(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := customerdto.CreateCustomerRequest{
		Name:     "João Silva",
		Email:    "joao@email.com",
		Password: "123456",
	}

	response := createCustomer(
		t,
		ts,
		token,
		request,
	)

	if response.ID == "" {
		t.Fatal("expected generated id")
	}

	if response.Name != request.Name {
		t.Fatal("invalid name")
	}

	if response.Email != request.Email {
		t.Fatal("invalid email")
	}
}

func TestCreateCustomerInvalidJSON(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
		token,
		bytes.NewBufferString(`{"name":`),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateCustomerInvalidName(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := customerdto.CreateCustomerRequest{
		Name:  "",
		Email: "cliente@email.com",
	}

	body := mustMarshal(t, request)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateCustomerInvalidEmail(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := customerdto.CreateCustomerRequest{
		Name:  "Cliente",
		Email: "",
	}

	body := mustMarshal(t, request)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateCustomerWithoutToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := customerdto.CreateCustomerRequest{
		Name:  "Cliente",
		Email: "cliente@email.com",
	}

	body := mustMarshal(t, request)

	req, err := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
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
func TestGetCustomerByID(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria Oliveira",
			Email: "maria@email.com",
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[customerdto.CustomerResponse](
		t,
		resp.Body,
	)

	if response.ID != customer.ID {
		t.Fatal("invalid id")
	}

	if response.Name != customer.Name {
		t.Fatal("invalid name")
	}

	if response.Email != customer.Email {
		t.Fatal("invalid email")
	}
}

func TestGetCustomerNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/customers/not-found",
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

func TestListCustomersEmpty(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/customers",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[[]customerdto.CustomerResponse](
		t,
		resp.Body,
	)

	if len(response) != 0 {
		t.Fatalf(
			"expected empty list got %d",
			len(response),
		)
	}
}

func TestListCustomers(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Maria",
			Email: "maria@email.com",
		},
	)

	createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "João",
			Email: "joao@email.com",
		},
	)

	createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Carlos",
			Email: "carlos@email.com",
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/customers",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[[]customerdto.CustomerResponse](
		t,
		resp.Body,
	)

	if len(response) != 3 {
		t.Fatalf(
			"expected 3 customers got %d",
			len(response),
		)
	}

	for _, customer := range response {
		if customer.ID == "" {
			t.Fatal("expected generated id")
		}

		if customer.Name == "" {
			t.Fatal("expected name")
		}

		if customer.Email == "" {
			t.Fatal("expected email")
		}
	}
}
func TestUpdateCustomer(t *testing.T) {

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

	update := customerdto.UpdateCustomerRequest{
		Name:  "João Silva",
		Email: "joao.silva@email.com",
	}

	body := mustMarshal(t, update)

	req, err := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[customerdto.CustomerResponse](
		t,
		resp.Body,
	)

	if response.Name != update.Name {
		t.Fatal("name not updated")
	}

	if response.Email != update.Email {
		t.Fatal("email not updated")
	}
}

func TestUpdateCustomerNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	update := customerdto.UpdateCustomerRequest{
		Name:  "Cliente",
		Email: "cliente@email.com",
	}

	body := mustMarshal(t, update)

	req, err := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/customers/not-found",
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestUpdateCustomerInvalidBody(t *testing.T) {

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

	req, err := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		bytes.NewBufferString(`{"name":`),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestDeleteCustomer(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Carlos",
			Email: "carlos@email.com",
		},
	)

	req, err := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNoContent(t, resp)

	req, err = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestDeleteCustomerNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/customers/not-found",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}
func TestCreateCustomerWithInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := customerdto.CreateCustomerRequest{
		Name:  "Cliente",
		Email: "cliente@email.com",
	}

	body := mustMarshal(t, request)

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/customers",
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

func TestUpdateCustomerInvalidValues(t *testing.T) {

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

	update := customerdto.UpdateCustomerRequest{
		Name:  "",
		Email: "",
	}

	body := mustMarshal(t, update)

	req, err := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestDeleteCustomerTwice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	customer := createCustomer(
		t,
		ts,
		token,
		customerdto.CreateCustomerRequest{
			Name:  "Carlos",
			Email: "carlos@email.com",
		},
	)

	req, _ := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/customers/"+customer.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestListCustomersUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/customers",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestGetCustomerUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/customers/123",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}
