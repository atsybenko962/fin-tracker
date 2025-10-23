package services

import (
	"github.com/atsybenko962/task-platform/user-service/internal/repository"
	g "github.com/task_platform/tools/user-service/grpc"
)

func grpcUserToCreateUserParams(request *g.CreateUserRequest) repository.CreateUserParams {
	return repository.CreateUserParams{
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	}
}

func UserToCreateUserResponse(user *repository.User) *g.CreateUserResponse {
	return &g.CreateUserResponse{
		UserId: user.ID.String(),
	}
}
