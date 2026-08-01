package memory

import (
	"testing"

	"desafio-go/internal/domain"
)

func newUserRepository() *MemoryUserRepository {
	return NewMemoryUserRepository().(*MemoryUserRepository)
}

func newValidUser() *domain.User {
	return domain.NewUser(
		"USR001",
		"John Doe",
		"john@email.com",
		"password_hash",
	)
}

func TestNewMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()

	if repo == nil {
		t.Fatal("expected repository")
	}
}

func TestUserSave(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	err := repo.Save(user)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(repo.users))
	}
}

func TestUserSave_InvalidUser(t *testing.T) {
	repo := newUserRepository()

	user := domain.NewUser(
		"",
		"",
		"",
		"",
	)

	err := repo.Save(user)

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(repo.users) != 0 {
		t.Fatal("invalid user should not be saved")
	}
}

func TestUserUpdate(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	_ = repo.Save(user)

	user.Name = "Jane Doe"

	err := repo.Update(user)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	saved, _ := repo.FindByID(user.ID)

	if saved.Name != "Jane Doe" {
		t.Fatal("user was not updated")
	}
}

func TestUserUpdate_NotFound(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	err := repo.Update(user)

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrUserNotFound,
			err,
		)
	}
}

func TestUserUpdate_InvalidUser(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	_ = repo.Save(user)

	user.Name = ""

	err := repo.Update(user)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUserDelete(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	_ = repo.Save(user)

	err := repo.Delete(user.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.users) != 0 {
		t.Fatal("user should have been removed")
	}
}

func TestUserDelete_NotFound(t *testing.T) {
	repo := newUserRepository()

	err := repo.Delete("INVALID")

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrUserNotFound,
			err,
		)
	}
}

func TestUserFindByID(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	_ = repo.Save(user)

	found, err := repo.FindByID(user.ID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != user.ID {
		t.Fatal("wrong user returned")
	}
}

func TestUserFindByID_NotFound(t *testing.T) {
	repo := newUserRepository()

	_, err := repo.FindByID("INVALID")

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrUserNotFound,
			err,
		)
	}
}

func TestUserFindByEmail(t *testing.T) {
	repo := newUserRepository()

	user := newValidUser()

	_ = repo.Save(user)

	found, err := repo.FindByEmail(user.Email)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != user.Email {
		t.Fatal("wrong user returned")
	}
}

func TestUserFindByEmail_NotFound(t *testing.T) {
	repo := newUserRepository()

	_, err := repo.FindByEmail("notfound@email.com")

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected %v, got %v",
			domain.ErrUserNotFound,
			err,
		)
	}
}

func TestUserFindAll(t *testing.T) {
	repo := newUserRepository()

	_ = repo.Save(domain.NewUser(
		"USR001",
		"John",
		"john@email.com",
		"hash1",
	))

	_ = repo.Save(domain.NewUser(
		"USR002",
		"Mary",
		"mary@email.com",
		"hash2",
	))

	users, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestUserFindAll_Empty(t *testing.T) {
	repo := newUserRepository()

	users, err := repo.FindAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 0 {
		t.Fatalf("expected empty slice, got %d", len(users))
	}
}
