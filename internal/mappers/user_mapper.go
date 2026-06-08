package mappers

import (
	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/entities"
)

func ToUserEntity(request dto.RegisterRequest) entities.User {
	return entities.User{
		Email:    request.Email,
		FullName: request.FullName,
	}
}

func ToUserResponse(user entities.User) dto.UserResponse {
	return dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}
}

func ToUpdateUserEntity(request dto.UpdateUserRequest) entities.User {
	return entities.User{
		FullName: request.FullName,
	}
}
