package auth

import (
	"context"
	"errors"
	"testing"
)

func TestLogoutUseCase_Execute(t *testing.T) {

	t.Run("should logout successfully", func(t *testing.T) {

		tokenStore := &tokenStoreStub{}

		useCase := NewLogoutUseCase(
			tokenStore,
		)

		err := useCase.Execute(
			context.Background(),
			"jwt-token",
		)

		if err != nil {
			t.Fatalf(
				"expected nil, got %v",
				err,
			)
		}
	})

	t.Run("should return error when token deletion fails", func(t *testing.T) {

		expectedErr := errors.New(
			"failed to delete token",
		)

		tokenStore := &tokenStoreStub{
			deleteError: expectedErr,
		}

		useCase := NewLogoutUseCase(
			tokenStore,
		)

		err := useCase.Execute(
			context.Background(),
			"jwt-token",
		)

		if !errors.Is(err, expectedErr) {
			t.Fatalf(
				"expected %v, got %v",
				expectedErr,
				err,
			)
		}
	})
}

func TestLogoutUseCase_Execute_ShouldDeleteCorrectToken(
	t *testing.T,
) {

	tokenStore := &tokenStoreSpy{}

	useCase := NewLogoutUseCase(
		tokenStore,
	)

	err := useCase.Execute(
		context.Background(),
		"jwt-token",
	)

	if err != nil {
		t.Fatalf(
			"expected nil, got %v",
			err,
		)
	}

	if !tokenStore.deleteCalled {
		t.Fatal("expected Delete to be called")
	}

	if tokenStore.deleteCalls != 1 {
		t.Fatalf(
			"expected 1 Delete call, got %d",
			tokenStore.deleteCalls,
		)
	}

	if tokenStore.token != "jwt-token" {
		t.Fatalf(
			"expected jwt-token, got %s",
			tokenStore.token,
		)
	}
}
