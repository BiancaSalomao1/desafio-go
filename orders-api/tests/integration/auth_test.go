/*
auth_helper

Responsabilidades:

- criar um usuário para testes;
- realizar login;
- obter um JWT válido;
- disponibilizar funções auxiliares para autenticação.

Este arquivo não contém testes.

Seu objetivo é evitar duplicação de código nos testes de integração.
*/

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	authdto "desafio-go/orders-api/internal/dto/auth"
	userdto "desafio-go/orders-api/internal/dto/user"
)

func createAuthenticatedUser(
	t *testing.T,
	ts *TestServer,
) string {

	t.Helper()

	email := "admin@test.com"
	password := "123456"

	createUserRequest := userdto.CreateUserRequest{
		Name:     "Administrator",
		Email:    email,
		Password: password,
	}

	body, err := json.Marshal(createUserRequest)
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}

	resp, err := http.Post(
		ts.Server.URL+"/users",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected create user status: %d", resp.StatusCode)
	}

	loginRequest := authdto.LoginRequest{
		Email:    email,
		Password: password,
	}

	body, err = json.Marshal(loginRequest)
	if err != nil {
		t.Fatalf("marshal login: %v", err)
	}

	resp, err = http.Post(
		ts.Server.URL+"/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected login status: %d", resp.StatusCode)
	}

	var response authdto.LoginResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("decode login: %v", err)
	}

	if response.AccessToken == "" {
		t.Fatal("empty access token")
	}

	return response.AccessToken
}
