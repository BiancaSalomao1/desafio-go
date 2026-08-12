package mapper

import (
	"testing"

	"desafio-go/orders-api/internal/domain"
	userdto "desafio-go/orders-api/internal/dto/user"
)

func TestToUser(t *testing.T) {

	request := userdto.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@email.com",
		Password: "123456",
	}

	user := ToUser(request)

	if user.ID == "" {
		t.Fatal("expected generated id")
	}

	if user.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, user.Name)
	}

	if user.Email != request.Email {
		t.Fatalf("expected %s, got %s", request.Email, user.Email)
	}

	if user.PasswordHash != request.Password {
		t.Fatalf("expected %s, got %s", request.Password, user.PasswordHash)
	}
}

func TestToUserResponse(t *testing.T) {

	user := domain.NewUser(
		"USR001",
		"John Doe",
		"john@email.com",
		"123456",
	)

	response := ToUserResponse(user)

	if response.ID != user.ID {
		t.Fatal("id mismatch")
	}

	if response.Name != user.Name {
		t.Fatal("name mismatch")
	}

	if response.Email != user.Email {
		t.Fatal("email mismatch")
	}
}

func TestToUserResponseList(t *testing.T) {

	users := []*domain.User{
		domain.NewUser(
			"USR001",
			"John Doe",
			"john@email.com",
			"123456",
		),
		domain.NewUser(
			"USR002",
			"Mary Doe",
			"mary@email.com",
			"654321",
		),
	}

	response := ToUserResponseList(users)

	if len(response) != 2 {
		t.Fatalf("expected 2 users, got %d", len(response))
	}

	if response[0].ID != "USR001" {
		t.Fatal("unexpected first user")
	}

	if response[1].ID != "USR002" {
		t.Fatal("unexpected second user")
	}
}

func TestToUserResponseList_Empty(t *testing.T) {

	response := ToUserResponseList([]*domain.User{})

	if len(response) != 0 {
		t.Fatalf("expected empty slice, got %d", len(response))
	}
}

func TestToUserUpdate(t *testing.T) {

	request := userdto.UpdateUserRequest{
		Name:  "Updated User",
		Email: "updated@email.com",
	}

	user := ToUserUpdate(
		"USR001",
		request,
	)

	if user.ID != "USR001" {
		t.Fatalf("expected USR001, got %s", user.ID)
	}

	if user.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, user.Name)
	}

	if user.Email != request.Email {
		t.Fatalf("expected %s, got %s", request.Email, user.Email)
	}

	if user.PasswordHash != "" {
		t.Fatalf("expected empty password hash, got %s", user.PasswordHash)
	}
}
