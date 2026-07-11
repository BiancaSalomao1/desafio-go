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
	ErrProductNotFound     = errors.New("product not found")
	ErrProductInvalid      = errors.New("product invalid")
	ErrProductNameRequired = errors.New("product name is required")
	ErrProductInUse        = errors.New("product is associated with one or more orders")

	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerInvalid  = errors.New("customer invalid")
	ErrCustomerInUse    = errors.New("customer is associated with one or more orders")

	ErrOrderNotFound       = errors.New("order not found")
	ErrEmptyOrder          = errors.New("empty order")
	ErrOrderStatusInvalid  = errors.New("order status invalid")
	ErrChangeStatusInvalid = errors.New("change status invalid")

	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInsufficientStock = errors.New("insufficient stock")

	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidStock = errors.New("invalid stock")

	ErrUserNotFound       = errors.New("user not found")
	ErrUserNameRequired   = errors.New("user name required")
	ErrUserInUse          = errors.New("user is associated with one or more records")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserEmailRequired  = errors.New("user email required")
	ErrPasswordRequired   = errors.New("password required")
)
