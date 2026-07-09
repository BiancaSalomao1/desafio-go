package domain

/*
Errors

Erros conhecidos:

- ErrProductNotFound
- ErrProductInvalid
- ErrCustomerNotFound
- ErrCustomerInvalid
- ErrUserNotFound
- ErrEmailAlreadyExists
- ErrOrderNotFound
- ErrEmptyOrder
- ErrOrderStatusInvalid
- ErrChangeStatusInvalid
- ErrInvalidQuantity
- ErrInsufficientStock
*/

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductInvalid  = errors.New("product invalid")

	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerInvalid  = errors.New("customer invalid")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrOrderNotFound       = errors.New("order not found")
	ErrEmptyOrder          = errors.New("empty order")
	ErrOrderStatusInvalid  = errors.New("order status invalid")
	ErrChangeStatusInvalid = errors.New("change status invalid")

	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInsufficientStock = errors.New("insufficient stock")

	ErrProductNameRequired = errors.New("product name is required")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidStock        = errors.New("invalid stock")

	ErrUserNameRequired  = errors.New("user name required")
	ErrUserEmailRequired = errors.New("user email required")
	ErrPasswordRequired  = errors.New("password required")
)
