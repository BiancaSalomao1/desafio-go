package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"orders-api/infrastructure/http/testutil"
	authdto "orders-api/internal/dto/auth"
	"testing"
)

// =========================
// Mock
// =========================

type loginUseCaseMock struct {
	token string
	err   error
}

func (m *loginUseCaseMock) Execute(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	return m.token, m.err
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

	t.Run("should return access token when credentials are valid", func(t *testing.T) {

		mock := &loginUseCaseMock{
			token: "jwt-token",
			err:   nil,
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
			token: "",
			err:   errors.New("invalid credentials"),
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
