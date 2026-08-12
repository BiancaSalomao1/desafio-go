package handler

import (
	"bytes"
	"desafio-go/orders-api/infrastructure/http/testutil"
	"desafio-go/orders-api/internal/domain"
	userdto "desafio-go/orders-api/internal/dto/user"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// =========================
// Mocks
// =========================

type createUserUseCaseMock struct {
	execute func(*domain.User) error
}

func (m *createUserUseCaseMock) Execute(user *domain.User) error {
	if m.execute != nil {
		return m.execute(user)
	}
	return nil
}

type getUserUseCaseMock struct {
	execute func(string) (*domain.User, error)
}

func (m *getUserUseCaseMock) Execute(id string) (*domain.User, error) {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil, nil
}

type listUsersUseCaseMock struct {
	execute func() ([]*domain.User, error)
}

func (m *listUsersUseCaseMock) Execute() ([]*domain.User, error) {
	if m.execute != nil {
		return m.execute()
	}
	return []*domain.User{}, nil
}

type updateUserUseCaseMock struct {
	execute func(*domain.User) error
}

func (m *updateUserUseCaseMock) Execute(user *domain.User) error {
	if m.execute != nil {
		return m.execute(user)
	}
	return nil
}

type deleteUserUseCaseMock struct {
	execute func(string) error
}

func (m *deleteUserUseCaseMock) Execute(id string) error {
	if m.execute != nil {
		return m.execute(id)
	}
	return nil
}

// =========================
// Helpers
// =========================

func newUserHandler(
	create CreateUserUseCase,
	get GetUserUseCase,
	list ListUsersUseCase,
	update UpdateUserUseCase,
	delete DeleteUserUseCase,
) *UserHandler {

	return NewUserHandler(
		create,
		get,
		list,
		update,
		delete,
	)
}

func validUser() *domain.User {
	return domain.NewUser(
		"USR001",
		"John Doe",
		"john@example.com",
		"123456",
	)
}
func TestUserHandler_Create(t *testing.T) {

	t.Run("should create user", func(t *testing.T) {

		mock := &createUserUseCaseMock{
			execute: func(user *domain.User) error {

				if user.Name != "John Doe" {
					t.Fatalf("unexpected name: %s", user.Name)
				}

				if user.Email != "john@example.com" {
					t.Fatalf("unexpected email: %s", user.Email)
				}

				return nil
			},
		}

		handler := newUserHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := userdto.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "123456",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/users",
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

		response := testutil.DecodeResponse[userdto.UserResponse](
			t,
			rec,
		)

		if response.Name != "John Doe" {
			t.Fatalf("expected John Doe, got %s", response.Name)
		}

		if response.Email != "john@example.com" {
			t.Fatalf("expected john@example.com, got %s", response.Email)
		}
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newUserHandler(
			&createUserUseCaseMock{},
			nil,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodPost,
			"/users",
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

		mock := &createUserUseCaseMock{
			execute: func(user *domain.User) error {
				return errors.New("user already exists")
			},
		}

		handler := newUserHandler(
			mock,
			nil,
			nil,
			nil,
			nil,
		)

		requestBody := userdto.CreateUserRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "123456",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPost,
			"/users",
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
			"user already exists",
		)
	})
}
func TestUserHandler_GetByID(t *testing.T) {

	t.Run("should return user", func(t *testing.T) {

		mock := &getUserUseCaseMock{
			execute: func(id string) (*domain.User, error) {

				if id != "USR001" {
					t.Fatalf("expected USR001, got %s", id)
				}

				return validUser(), nil
			},
		}

		handler := newUserHandler(
			nil,
			mock,
			nil,
			nil,
			nil,
		)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/USR001",
			nil,
		)

		req.SetPathValue("id", "USR001")

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		testutil.AssertStatus(t, rec, http.StatusOK)
		testutil.AssertContentTypeJSON(t, rec)

		response := testutil.DecodeResponse[userdto.UserResponse](t, rec)

		if response.ID != "USR001" {
			t.Fatalf("expected USR001, got %s", response.ID)
		}
	})

	t.Run("should return not found", func(t *testing.T) {

		mock := &getUserUseCaseMock{
			execute: func(id string) (*domain.User, error) {
				return nil, errors.New("user not found")
			},
		}

		handler := newUserHandler(nil, mock, nil, nil, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/USR001",
			nil,
		)

		req.SetPathValue("id", "USR001")

		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		testutil.AssertStatus(t, rec, http.StatusNotFound)
		testutil.AssertContentTypeJSON(t, rec)
		testutil.AssertBodyContains(t, rec, "user not found")
	})
}
func TestUserHandler_List(t *testing.T) {

	t.Run("should list users", func(t *testing.T) {

		mock := &listUsersUseCaseMock{
			execute: func() ([]*domain.User, error) {
				return []*domain.User{
					validUser(),
				}, nil
			},
		}

		handler := newUserHandler(nil, nil, mock, nil, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(rec, req)

		testutil.AssertStatus(t, rec, http.StatusOK)
		testutil.AssertContentTypeJSON(t, rec)

		response := testutil.DecodeResponse[[]userdto.UserResponse](t, rec)

		if len(response) != 1 {
			t.Fatalf("expected 1 user, got %d", len(response))
		}
	})

	t.Run("should return internal server error", func(t *testing.T) {

		mock := &listUsersUseCaseMock{
			execute: func() ([]*domain.User, error) {
				return nil, errors.New("database error")
			},
		}

		handler := newUserHandler(nil, nil, mock, nil, nil)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users",
			nil,
		)

		rec := httptest.NewRecorder()

		handler.List(rec, req)

		testutil.AssertStatus(t, rec, http.StatusInternalServerError)
		testutil.AssertContentTypeJSON(t, rec)
		testutil.AssertBodyContains(t, rec, "database error")
	})
}
func TestUserHandler_Update(t *testing.T) {

	t.Run("should update user", func(t *testing.T) {

		mock := &updateUserUseCaseMock{
			execute: func(user *domain.User) error {

				if user.ID != "USR001" {
					t.Fatalf("expected USR001")
				}

				return nil
			},
		}

		handler := newUserHandler(nil, nil, nil, mock, nil)

		body := userdto.UpdateUserRequest{
			Name:  "John Updated",
			Email: "john@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/users/USR001",
			body,
		)

		req.SetPathValue("id", "USR001")

		handler.Update(rec, req)

		testutil.AssertStatus(t, rec, http.StatusOK)
		testutil.AssertContentTypeJSON(t, rec)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {

		handler := newUserHandler(nil, nil, nil, &updateUserUseCaseMock{}, nil)

		req := httptest.NewRequest(
			http.MethodPut,
			"/users/USR001",
			bytes.NewBufferString("{invalid"),
		)

		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", "USR001")

		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		testutil.AssertStatus(t, rec, http.StatusBadRequest)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &updateUserUseCaseMock{
			execute: func(user *domain.User) error {
				return errors.New("cannot update")
			},
		}

		handler := newUserHandler(nil, nil, nil, mock, nil)

		body := userdto.UpdateUserRequest{
			Name:  "John",
			Email: "john@example.com",
		}

		req, rec := testutil.NewJSONRequest(
			http.MethodPut,
			"/users/USR001",
			body,
		)

		req.SetPathValue("id", "USR001")

		handler.Update(rec, req)

		testutil.AssertStatus(t, rec, http.StatusBadRequest)
		testutil.AssertBodyContains(t, rec, "cannot update")
	})
}
func TestUserHandler_Delete(t *testing.T) {

	t.Run("should delete user", func(t *testing.T) {

		mock := &deleteUserUseCaseMock{
			execute: func(id string) error {

				if id != "USR001" {
					t.Fatalf("expected USR001")
				}

				return nil
			},
		}

		handler := newUserHandler(nil, nil, nil, nil, mock)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/users/USR001",
			nil,
		)

		req.SetPathValue("id", "USR001")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		testutil.AssertStatus(t, rec, http.StatusNoContent)
	})

	t.Run("should return bad request when usecase returns error", func(t *testing.T) {

		mock := &deleteUserUseCaseMock{
			execute: func(id string) error {
				return errors.New("cannot delete")
			},
		}

		handler := newUserHandler(nil, nil, nil, nil, mock)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/users/USR001",
			nil,
		)

		req.SetPathValue("id", "USR001")

		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		testutil.AssertStatus(t, rec, http.StatusBadRequest)
		testutil.AssertBodyContains(t, rec, "cannot delete")
	})
}
