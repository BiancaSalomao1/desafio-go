/*
product_test

Responsabilidades:

- testar toda a API de Produtos;
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

Cenários cobertos nesta parte:

✓ Create Product
✓ Create Invalid JSON
✓ Create Invalid Name
✓ Create Invalid Price
✓ Create Invalid Stock
*/

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	productdto "desafio-go/internal/dto/product"
)

func createProduct(
	t *testing.T,
	ts *TestServer,
	token string,
	request productdto.CreateProductRequest,
) productdto.ProductResponse {

	t.Helper()

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		token,
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatalf("request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {

		t.Fatalf(
			"expected %d got %d",
			http.StatusCreated,
			resp.StatusCode,
		)

	}

	var response productdto.ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	return response
}

func TestCreateProduct(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := productdto.CreateProductRequest{
		Name:  "Notebook Dell",
		Price: 5200,
		Stock: 8,
	}

	response := createProduct(
		t,
		ts,
		token,
		request,
	)

	if response.ID == "" {
		t.Fatal("expected id")
	}

	if response.Name != request.Name {
		t.Fatal("invalid name")
	}

	if response.Price != request.Price {
		t.Fatal("invalid price")
	}

	if response.Stock != request.Stock {
		t.Fatal("invalid stock")
	}

}

func TestCreateProductInvalidJSON(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		token,
		bytes.NewBufferString(`{"name":`),
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}

}

func TestCreateProductInvalidName(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := productdto.CreateProductRequest{
		Name:  "",
		Price: 100,
		Stock: 10,
	}

	body, _ := json.Marshal(request)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		token,
		bytes.NewBuffer(body),
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}

}

func TestCreateProductInvalidPrice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := productdto.CreateProductRequest{
		Name:  "Mouse",
		Price: 0,
		Stock: 5,
	}

	body, _ := json.Marshal(request)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		token,
		bytes.NewBuffer(body),
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}

}

func TestCreateProductInvalidStock(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	request := productdto.CreateProductRequest{
		Name:  "Monitor",
		Price: 1200,
		Stock: -1,
	}

	body, _ := json.Marshal(request)

	req, _ := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		token,
		bytes.NewBuffer(body),
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}

}
func TestGetProductByID(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "MacBook Pro",
			Price: 15000,
			Stock: 3,
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	var response productdto.ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.ID != product.ID {
		t.Fatal("invalid id")
	}

	if response.Name != product.Name {
		t.Fatal("invalid name")
	}

	if response.Price != product.Price {
		t.Fatal("invalid price")
	}

	if response.Stock != product.Stock {
		t.Fatal("invalid stock")
	}

}

func TestGetProductNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/not-found",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {

		t.Fatalf(
			"expected %d got %d",
			http.StatusNotFound,
			resp.StatusCode,
		)

	}

}

func TestListProductsEmpty(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			resp.StatusCode,
		)

	}

	var response []productdto.ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response) != 0 {
		t.Fatalf(
			"expected empty list got %d",
			len(response),
		)
	}

}

func TestListProducts(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 5000,
			Stock: 5,
		},
	)

	createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 100,
			Stock: 30,
		},
	)

	createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Teclado",
			Price: 200,
			Stock: 15,
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products",
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		t.Fatalf(
			"expected %d got %d",
			http.StatusOK,
			resp.StatusCode,
		)

	}

	var response []productdto.ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if len(response) != 3 {

		t.Fatalf(
			"expected 3 products got %d",
			len(response),
		)

	}

	if response[0].ID == "" {
		t.Fatal("expected generated id")
	}

}
func TestUpdateProduct(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Notebook",
			Price: 4500,
			Stock: 5,
		},
	)

	update := productdto.UpdateProductRequest{
		Name:  "Notebook Gamer",
		Price: 6500,
		Stock: 8,
	}

	body, err := json.Marshal(update)
	if err != nil {
		t.Fatal(err)
	}

	req, err := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/products/"+product.ID,
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

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d",
			http.StatusOK,
			resp.StatusCode,
		)
	}

	var response productdto.ProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Name != update.Name {
		t.Fatal("name not updated")
	}

	if response.Price != update.Price {
		t.Fatal("price not updated")
	}

	if response.Stock != update.Stock {
		t.Fatal("stock not updated")
	}
}

func TestUpdateProductNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	update := productdto.UpdateProductRequest{
		Name:  "Notebook",
		Price: 5000,
		Stock: 5,
	}

	body, _ := json.Marshal(update)

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/products/not-found",
		token,
		bytes.NewBuffer(body),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)
	}
}

func TestUpdateProductInvalidBody(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

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

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/products/"+product.ID,
		token,
		bytes.NewBufferString(`{"name":`),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)
	}
}

func TestDeleteProduct(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 100,
			Stock: 20,
		},
	)

	req, err := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected %d got %d",
			http.StatusNoContent,
			resp.StatusCode,
		)
	}

	// garante que foi removido
	req, _ = authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected %d got %d",
			http.StatusNotFound,
			resp.StatusCode,
		)
	}
}

func TestDeleteProductNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, _ := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/products/not-found",
		token,
		nil,
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)
	}
}
func TestCreateProductWithoutToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := productdto.CreateProductRequest{
		Name:  "Notebook",
		Price: 5000,
		Stock: 10,
	}

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			resp.StatusCode,
		)
	}
}

func TestCreateProductWithInvalidToken(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := productdto.CreateProductRequest{
		Name:  "Notebook",
		Price: 5000,
		Stock: 10,
	}

	body, _ := json.Marshal(request)

	req, err := authenticatedRequest(
		http.MethodPost,
		ts.Server.URL+"/products",
		"invalid-token",
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

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf(
			"expected %d got %d",
			http.StatusUnauthorized,
			resp.StatusCode,
		)
	}
}

func TestUpdateProductInvalidValues(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Mouse",
			Price: 100,
			Stock: 10,
		},
	)

	update := productdto.UpdateProductRequest{
		Name:  "",
		Price: -1,
		Stock: -5,
	}

	body, _ := json.Marshal(update)

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/products/"+product.ID,
		token,
		bytes.NewBuffer(body),
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}
}

func TestDeleteProductTwice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	product := createProduct(
		t,
		ts,
		token,
		productdto.CreateProductRequest{
			Name:  "Keyboard",
			Price: 250,
			Stock: 30,
		},
	)

	req, _ := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf(
			"expected %d got %d",
			http.StatusNoContent,
			resp.StatusCode,
		)
	}

	req, _ = authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/products/"+product.ID,
		token,
		nil,
	)

	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {

		t.Fatalf(
			"expected %d got %d",
			http.StatusBadRequest,
			resp.StatusCode,
		)

	}
}
