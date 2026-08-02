package memory

import (
	"testing"

	"desafio-go/internal/domain"
)

func newCustomerRepository() *MemoryCustomerRepository {
	return NewMemoryCustomerRepository().(*MemoryCustomerRepository)
}

func newValidCustomer() *domain.Customer {
	return domain.NewCustomer(
		"CUS001",
		"João Silva",
		"joao@email.com",
		"joaohash",
	)
}

func TestNewMemoryCustomerRepository(t *testing.T) {
	repo := NewMemoryCustomerRepository()

	if repo == nil {
		t.Fatal("expected repository")
	}
}

func TestCustomerSave(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	err := repo.Save(customer)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.customers) != 1 {
		t.Fatalf("expected 1 customer, got %d", len(repo.customers))
	}
}

func TestCustomerSave_InvalidCustomer(t *testing.T) {
	repo := newCustomerRepository()

	customer := domain.NewCustomer(
		"",
		"",
		"",
		"",
	)

	err := repo.Save(customer)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(repo.customers) != 0 {
		t.Fatal("invalid customer should not be saved")
	}
}

func TestCustomerUpdate(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	_ = repo.Save(customer)

	customer.Name = "Maria Silva"

	err := repo.Update(customer)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, _ := repo.FindByID(customer.ID)

	if saved.Name != "Maria Silva" {
		t.Fatal("customer was not updated")
	}
}

func TestCustomerUpdate_NotFound(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	err := repo.Update(customer)

	if err != domain.ErrCustomerNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrCustomerNotFound,
			err,
		)
	}
}

func TestCustomerUpdate_InvalidCustomer(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	_ = repo.Save(customer)

	customer.Name = ""

	err := repo.Update(customer)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCustomerDelete(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	_ = repo.Save(customer)

	err := repo.Delete(customer.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.customers) != 0 {
		t.Fatal("customer should have been removed")
	}
}

func TestCustomerDelete_NotFound(t *testing.T) {
	repo := newCustomerRepository()

	err := repo.Delete("INVALID")

	if err != domain.ErrCustomerNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrCustomerNotFound,
			err,
		)
	}
}

func TestCustomerFindByID(t *testing.T) {
	repo := newCustomerRepository()

	customer := newValidCustomer()

	_ = repo.Save(customer)

	found, err := repo.FindByID(customer.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != customer.ID {
		t.Fatal("wrong customer returned")
	}
}

func TestCustomerFindByID_NotFound(t *testing.T) {
	repo := newCustomerRepository()

	_, err := repo.FindByID("INVALID")

	if err != domain.ErrCustomerNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrCustomerNotFound,
			err,
		)
	}
}

func TestCustomerFindAll(t *testing.T) {
	repo := newCustomerRepository()

	_ = repo.Save(domain.NewCustomer(
		"1",
		"João",
		"joao@email.com",
		"joaohash",
	))

	_ = repo.Save(domain.NewCustomer(
		"2",
		"Maria",
		"maria@email.com",
		"mariahash",
	))

	customers, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(customers) != 2 {
		t.Fatalf("expected 2 customers, got %d", len(customers))
	}
}

func TestCustomerFindAll_Empty(t *testing.T) {
	repo := newCustomerRepository()

	customers, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(customers) != 0 {
		t.Fatalf("expected empty slice, got %d", len(customers))
	}
}
