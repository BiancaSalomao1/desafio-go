package httpx

import (
	"errors"

	"desafio-go/internal/domain"
)

func errorMessage(err error) string {

	switch {

	case errors.Is(err, domain.ErrCustomerNotFound):
		return err.Error()

	case errors.Is(err, domain.ErrProductNotFound):
		return err.Error()

	case errors.Is(err, domain.ErrOrderNotFound):
		return err.Error()

	case errors.Is(err, domain.ErrInsufficientStock):
		return err.Error()

	case errors.Is(err, domain.ErrOrderStatusInvalid):
		return err.Error()

	case errors.Is(err, domain.ErrChangeStatusInvalid):
		return err.Error()

	case errors.Is(err, domain.ErrDuplicatedProduct):
		return err.Error()

	case errors.Is(err, domain.ErrCustomerInvalid):
		return err.Error()

	case errors.Is(err, domain.ErrProductInvalid):
		return err.Error()

	case errors.Is(err, domain.ErrEmptyOrder):
		return err.Error()

	case errors.Is(err, domain.ErrInvalidQuantity):
		return err.Error()

	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return err.Error()

	case errors.Is(err, domain.ErrInvalidCredentials):
		return err.Error()

	default:
		return "internal server error"
	}
}
