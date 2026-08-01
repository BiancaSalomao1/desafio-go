package mapper

import (
	"testing"

	"desafio-go/internal/domain"
	customerdto "desafio-go/internal/dto/customer"
)

func TestToCustomer(t *testing.T) {

	request := customerdto.CreateCustomerRequest{
		Name:  "João Silva",
		Email: "joao@email.com",
	}

	customer := ToCustomer(request)

	if customer.ID == "" {
		t.Fatal("expected generated id")
	}

	if customer.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, customer.Name)
	}

	if customer.Email != request.Email {
		t.Fatalf("expected %s, got %s", request.Email, customer.Email)
	}
}

func TestToCustomerUpdate(t *testing.T) {

	request := customerdto.UpdateCustomerRequest{
		Name:  "Maria Souza",
		Email: "maria@email.com",
	}

	customer := ToCustomerUpdate(
		"CUST001",
		request,
	)

	if customer.ID != "CUST001" {
		t.Fatalf("expected CUST001, got %s", customer.ID)
	}

	if customer.Name != request.Name {
		t.Fatalf("expected %s, got %s", request.Name, customer.Name)
	}

	if customer.Email != request.Email {
		t.Fatalf("expected %s, got %s", request.Email, customer.Email)
	}
}

func TestToCustomerResponse(t *testing.T) {

	customer := domain.NewCustomer(
		"CUST001",
		"João Silva",
		"joao@email.com",
	)

	response := ToCustomerResponse(customer)

	if response.ID != customer.ID {
		t.Fatal("id mismatch")
	}

	if response.Name != customer.Name {
		t.Fatal("name mismatch")
	}

	if response.Email != customer.Email {
		t.Fatal("email mismatch")
	}
}

func TestToCustomerResponseList(t *testing.T) {

	customers := []*domain.Customer{
		domain.NewCustomer(
			"CUST001",
			"João Silva",
			"joao@email.com",
		),
		domain.NewCustomer(
			"CUST002",
			"Maria Souza",
			"maria@email.com",
		),
	}

	response := ToCustomerResponseList(customers)

	if len(response) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(response))
	}

	if response[0].ID != "CUST001" {
		t.Fatal("unexpected first customer")
	}

	if response[1].ID != "CUST002" {
		t.Fatal("unexpected second customer")
	}
}

func TestToCustomerResponseList_Empty(t *testing.T) {

	response := ToCustomerResponseList([]*domain.Customer{})

	if len(response) != 0 {
		t.Fatalf("expected empty slice, got %d", len(response))
	}
}
