package domain

/*
Errors

Erros Products conhecidos:

- ErrProductNotFound
- ErrProductInvalid
- ErrProductNameRequired
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

	ErrOrderNotFound       = errors.New("order not found")
	ErrEmptyOrder          = errors.New("empty order")
	ErrDuplicatedProduct   = errors.New("duplicated product in order")
	ErrOrderStatusInvalid  = errors.New("order status invalid")
	ErrChangeStatusInvalid = errors.New("change status invalid")

	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrInsufficientStock = errors.New("insufficient stock")

	ErrInvalidPrice = errors.New("invalid price")
	ErrInvalidStock = errors.New("invalid stock")

	ErrInternal = errors.New("internal server error")
)
