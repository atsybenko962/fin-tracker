package services

import (
	"fmt"
	"github.com/atsybenko962/task-platform/user-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	g "github.com/task_platform/tools/user-service/grpc"
	"time"
)

func grpcUserToCreateUserParams(request *g.CreateUserRequest) repository.CreateUserParams {
	return repository.CreateUserParams{
		Email:        request.Email,
		PasswordHash: request.PasswordHash,
	}
}

func grpcUuidToUUID(request *g.DeleteUserRequest) (pgtype.UUID, error) {
	stdUUID, err := uuid.Parse(request.UserId)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}

	return pgtype.UUID{
		Bytes: [16]byte(stdUUID),
		Valid: true,
	}, nil
}

func UserToCreateUserResponse(user *repository.User) *g.CreateUserResponse {
	return &g.CreateUserResponse{
		UserId: user.ID.String(),
	}
}

func UserToGetUserResponse(user *repository.User) *g.GetUserResponse {
	return &g.GetUserResponse{
		UserId:       user.ID.String(),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.Time.Format(time.RFC3339Nano),
	}
}
