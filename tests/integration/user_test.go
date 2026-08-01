/*
user_test

Responsabilidades:

- testar toda a API de Usuários;
- validar autenticação JWT;
- validar integração entre HTTP e PostgreSQL.

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

✓ Create User
✓ Create Invalid JSON
✓ Create Invalid Name
✓ Create Invalid Email
✓ Create Invalid Password
✓ Duplicate Email
*/

package integration

import (
	"bytes"
	"net/http"
	"testing"

	userdto "desafio-go/internal/dto/user"
)

func createUser(
	t *testing.T,
	ts *TestServer,
	request userdto.CreateUserRequest,
) userdto.UserResponse {

	t.Helper()

	body := mustMarshal(t, request)

	req, err := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBuffer(body),
	)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertCreated(t, resp)

	return decodeResponse[userdto.UserResponse](
		t,
		resp.Body,
	)
}

func TestCreateUser(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := userdto.CreateUserRequest{
		Name:     "Administrador",
		Email:    "admin@email.com",
		Password: "123456",
	}

	response := createUser(t, ts, request)

	if response.ID == "" {
		t.Fatal("expected id")
	}

	if response.Name != request.Name {
		t.Fatal("invalid name")
	}

	if response.Email != request.Email {
		t.Fatal("invalid email")
	}
}

func TestCreateUserInvalidJSON(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, _ := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBufferString(`{"name":`),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateUserInvalidName(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := userdto.CreateUserRequest{
		Name:     "",
		Email:    "admin@email.com",
		Password: "123456",
	}

	body := mustMarshal(t, request)

	req, _ := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateUserInvalidEmail(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := userdto.CreateUserRequest{
		Name:     "Administrador",
		Email:    "",
		Password: "123456",
	}

	body := mustMarshal(t, request)

	req, _ := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateUserInvalidPassword(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := userdto.CreateUserRequest{
		Name:     "Administrador",
		Email:    "admin@email.com",
		Password: "",
	}

	body := mustMarshal(t, request)

	req, _ := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestCreateUserDuplicateEmail(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	request := userdto.CreateUserRequest{
		Name:     "Administrador",
		Email:    "admin@email.com",
		Password: "123456",
	}

	createUser(t, ts, request)

	body := mustMarshal(t, request)

	req, _ := http.NewRequest(
		http.MethodPost,
		ts.Server.URL+"/users",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}
func TestGetUserByID(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	user := createUser(
		t,
		ts,
		userdto.CreateUserRequest{
			Name:     "Maria",
			Email:    "maria@email.com",
			Password: "123456",
		},
	)

	req, err := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/users/"+user.ID,
		token,
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[userdto.UserResponse](t, resp.Body)

	if response.ID != user.ID {
		t.Fatal("invalid id")
	}

	if response.Name != user.Name {
		t.Fatal("invalid name")
	}

	if response.Email != user.Email {
		t.Fatal("invalid email")
	}
}

func TestGetUserNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, _ := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/users/not-found",
		token,
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertNotFound(t, resp)
}

func TestListUsers(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	createUser(t, ts, userdto.CreateUserRequest{
		Name:     "João",
		Email:    "joao@email.com",
		Password: "123456",
	})

	createUser(t, ts, userdto.CreateUserRequest{
		Name:     "Maria",
		Email:    "maria@email.com",
		Password: "123456",
	})

	req, _ := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/users",
		token,
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	users := decodeResponse[[]userdto.UserResponse](t, resp.Body)

	if len(users) != 3 {
		t.Fatalf("expected 3 users got %d", len(users))
	}
}

func TestListUsersEmpty(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, _ := authenticatedRequest(
		http.MethodGet,
		ts.Server.URL+"/users",
		token,
		nil,
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	users := decodeResponse[[]userdto.UserResponse](t, resp.Body)

	if len(users) != 1 {
		t.Fatalf("expected 1 user got %d", len(users))
	}
}

func TestUpdateUser(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	user := createUser(
		t,
		ts,
		userdto.CreateUserRequest{
			Name:     "João",
			Email:    "joao@email.com",
			Password: "123456",
		},
	)

	update := userdto.UpdateUserRequest{
		Name:  "João Silva",
		Email: "joao.silva@email.com",
	}

	body := mustMarshal(t, update)

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/users/"+user.ID,
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertOK(t, resp)

	response := decodeResponse[userdto.UserResponse](t, resp.Body)

	if response.Name != update.Name {
		t.Fatal("name not updated")
	}

	if response.Email != update.Email {
		t.Fatal("email not updated")
	}
}

func TestUpdateUserNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	update := userdto.UpdateUserRequest{
		Name:  "Novo Nome",
		Email: "novo@email.com",
	}

	body := mustMarshal(t, update)

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/users/not-found",
		token,
		bytes.NewBuffer(body),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestUpdateUserInvalidBody(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	user := createUser(
		t,
		ts,
		userdto.CreateUserRequest{
			Name:     "Maria",
			Email:    "maria@email.com",
			Password: "123456",
		},
	)

	req, _ := authenticatedRequest(
		http.MethodPut,
		ts.Server.URL+"/users/"+user.ID,
		token,
		bytes.NewBufferString(`{"name":`),
	)

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}
func TestDeleteUser(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	user := createUser(
		t,
		ts,
		userdto.CreateUserRequest{
			Name:     "Carlos",
			Email:    "carlos@email.com",
			Password: "123456",
		},
	)

	req, err := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/users/"+user.ID,
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
		ts.Server.URL+"/users/"+user.ID,
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

func TestDeleteUserNotFound(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	req, err := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/users/not-found",
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

func TestDeleteUserTwice(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	token := createAuthenticatedUser(t, ts)

	user := createUser(
		t,
		ts,
		userdto.CreateUserRequest{
			Name:     "Maria",
			Email:    "maria@email.com",
			Password: "123456",
		},
	)

	req, _ := authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/users/"+user.ID,
		token,
		nil,
	)

	resp := doRequest(t, req)
	resp.Body.Close()

	assertNoContent(t, resp)

	req, _ = authenticatedRequest(
		http.MethodDelete,
		ts.Server.URL+"/users/"+user.ID,
		token,
		nil,
	)

	resp = doRequest(t, req)
	defer resp.Body.Close()

	assertBadRequest(t, resp)
}

func TestListUsersUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/users",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestGetUserUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodGet,
		ts.Server.URL+"/users/123",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestUpdateUserUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodPut,
		ts.Server.URL+"/users/123",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}

func TestDeleteUserUnauthorized(t *testing.T) {

	ts := setup(t)
	defer teardown(ts)

	req, err := http.NewRequest(
		http.MethodDelete,
		ts.Server.URL+"/users/123",
		nil,
	)

	if err != nil {
		t.Fatal(err)
	}

	resp := doRequest(t, req)
	defer resp.Body.Close()

	assertUnauthorized(t, resp)
}
