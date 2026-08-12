package handler

import (
	"bytes"
	"desafio-go/orders-api/infrastructure/http/testutil"
	authdto "desafio-go/orders-api/internal/dto/auth"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =========================
// Mock
// =========================

type loginUseCaseMock struct {
	execute func(email, password string) (string, error)
}

func (m *loginUseCaseMock) Execute(
	email,
	password string,
) (string, error) {

	if m.execute != nil {
		return m.execute(email, password)
	}

	return "", nil
}

// =========================
// Helper
// =========================

func newAuthHandler(
	login LoginUseCase,
) *AuthHandler {

	return NewAuthHandler(login)
}

// =========================
// Tests
// =========================

func TestAuthHandler_Login(t *testing.T) {

	t.Run("should login successfully", func(t *testing.T) {

		mock := &loginUseCaseMock{
			execute: func(email, password string) (string, error) {

				if email != "john@example.com" {
					t.Fatalf("unexpected email: %s", email)
				}

				if password != "123456" {
					t.Fatalf("unexpected password")
				}

				return "jwt-token", nil
			},
		}

		handler := newAuthHandler(mock)

		requestBody := authdto.LoginRequest{
			Email:    "john@example.com",
			Password: "123456",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/login",
			requestBody,
		)

		handler.Login(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusOK,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		response := testutil.DecodeResponse[authdto.LoginResponse](
			t,
			rec,
		)

		if response.AccessToken != "jwt-token" {
			t.Fatalf(
				"expected jwt-token, got %s",
				response.AccessToken,
			)
		}
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newAuthHandler(&loginUseCaseMock{})

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBufferString("{invalid"),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

		rec := httptest.NewRecorder()

		handler.Login(rec, req)

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

	t.Run("should return unauthorized", func(t *testing.T) {

		mock := &loginUseCaseMock{
			execute: func(email, password string) (string, error) {
				return "", errors.New("invalid credentials")
			},
		}

		handler := newAuthHandler(mock)

		requestBody := authdto.LoginRequest{
			Email:    "john@example.com",
			Password: "wrong-password",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/login",
			requestBody,
		)

		handler.Login(rec, req)

		testutil.AssertStatus(
			t,
			rec,
			http.StatusUnauthorized,
		)

		testutil.AssertContentTypeJSON(
			t,
			rec,
		)

		testutil.AssertBodyContains(
			t,
			rec,
			"invalid credentials",
		)
	})
}
