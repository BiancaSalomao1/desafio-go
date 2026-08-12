package mapper

import (
	"desafio-go/orders-api/internal/domain"
	userdto "desafio-go/orders-api/internal/dto/user"

	"github.com/google/uuid"
)

func ToUser(
	request userdto.CreateUserRequest,
) *domain.User {

	return domain.NewUser(
		uuid.NewString(),
		request.Name,
		request.Email,
		request.Password,
	)
}

func ToUserResponse(
	user *domain.User,
) userdto.UserResponse {

	return userdto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func ToUserResponseList(
	users []*domain.User,
) []userdto.UserResponse {

	response := make(
		[]userdto.UserResponse,
		0,
		len(users),
	)

	for _, user := range users {

		response = append(
			response,
			ToUserResponse(user),
		)
	}

	return response
}

func ToUserUpdate(
	id string,
	request userdto.UpdateUserRequest,
) *domain.User {

	return domain.NewUser(
		id,
		request.Name,
		request.Email,
		"",
	)
}
