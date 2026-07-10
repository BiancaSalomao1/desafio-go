package mapper

import (
	"desafio-go/internal/domain"
	customerdto "desafio-go/internal/dto/customer"

	"github.com/google/uuid"
)

func ToCustomer(
	request customerdto.CreateCustomerRequest,
) *domain.Customer {

	return domain.NewCustomer(
		uuid.NewString(),
		request.Name,
		request.Email,
	)
}

func ToUpdatedCustomer(
	id string,
	request customerdto.UpdateCustomerRequest,
) *domain.Customer {

	return domain.NewCustomer(
		id,
		request.Name,
		request.Email,
	)
}

func ToCustomerResponse(
	customer *domain.Customer,
) customerdto.CustomerResponse {

	return customerdto.CustomerResponse{
		ID:    customer.ID,
		Name:  customer.Name,
		Email: customer.Email,
	}
}

func ToCustomerResponseList(
	customers []*domain.Customer,
) []customerdto.CustomerResponse {

	response := make(
		[]customerdto.CustomerResponse,
		0,
		len(customers),
	)

	for _, customer := range customers {

		response = append(
			response,
			ToCustomerResponse(customer),
		)
	}

	return response
}
