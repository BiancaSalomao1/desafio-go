package mapper

/*
Função ToCustomer

Responsabilidades:
- converter DTOs de cliente para domínio;
- converter domínio para DTO de resposta.

Métodos:
- ToCustomer()
- ToUpdatedCustomer()
- ToCustomerResponse()
- ToCustomerResponseList()
- ToCustomerUpdate()
*/

import (
	"orders-api/internal/domain"
	customerdto "orders-api/internal/dto/customer"

	"github.com/google/uuid"
)

func ToCustomer(
	request customerdto.CreateCustomerRequest,
) *domain.Customer {

	customer := domain.NewCustomer(
		uuid.NewString(),
		request.Name,
		request.Email,
		"",
	)

	customer.Password = request.Password

	return customer
}

func ToUpdatedCustomer(
	id string,
	request customerdto.UpdateCustomerRequest,
) *domain.Customer {

	return domain.NewCustomer(
		id,
		request.Name,
		request.Email,
		"",
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

func ToCustomerUpdate(
	id string,
	request customerdto.UpdateCustomerRequest,
) *domain.Customer {

	return domain.NewCustomer(
		id,
		request.Name,
		request.Email,
		"",
	)
}
