package services

import (
	"context"
	"github.com/atsybenko962/task-platform/user-service/internal/repository"
	errors "github.com/task_platform/tools/helpers"
	"github.com/task_platform/tools/httptools"
	g "github.com/task_platform/tools/user-service/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const DefaultServiceTimeout = 3 * time.Second

type UserService struct {
	g.UnimplementedUserServiceServer
	repo   *repository.Queries
	logger *zap.Logger
}

func NewUserService(db repository.DBTX, logger *zap.Logger) *UserService {
	return &UserService{
		UnimplementedUserServiceServer: g.UnimplementedUserServiceServer{},
		repo:                           repository.New(db),
		logger:                         logger,
	}
}

func (s *UserService) CreateUserGrpc(ctx context.Context, request *g.CreateUserRequest) (*g.CreateUserResponse, error) {
	info := "grpc user creation service"

	if request == nil {
		s.logger.Warn(info, zap.String("error", errors.StatusInvalidArgumentError), zap.Any("request", request))
		return nil, status.Error(codes.InvalidArgument, errors.StatusInvalidArgumentError)
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultServiceTimeout)
	defer cancel()

	ch := make(chan *repository.User)
	chErr := make(chan error)

	go func() {
		defer close(ch)
		defer close(chErr)

		result, err := s.repo.CreateUser(ctx, grpcUserToCreateUserParams(request))
		if err != nil {
			chErr <- err
		}

		ch <- &result
		return
	}()

	user, err := httptools.SelectChannels[repository.User](ctx, ch, chErr, info, request, s.logger)
	return UserToCreateUserResponse(user), err

}

func (s *UserService) GetUserByEmailGrpc(ctx context.Context, request *g.GetUserRequest) (*g.GetUserResponse, error) {
	info := "grpc service for searching for a user by email"

	if request == nil {
		s.logger.Warn(info, zap.String("error", errors.StatusInvalidArgumentError), zap.Any("request", request))
		return nil, status.Error(codes.InvalidArgument, errors.StatusInvalidArgumentError)
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultServiceTimeout)
	defer cancel()

	ch := make(chan *repository.User)
	chErr := make(chan error)

	go func() {
		defer close(ch)
		defer close(chErr)

		result, err := s.repo.GetUserByEmail(ctx, request.String())
		if err != nil {
			chErr <- err
		}

		ch <- &result
		return
	}()

	user, err := httptools.SelectChannels[repository.User](ctx, ch, chErr, info, request, s.logger)
	if err != nil {
		return nil, err
	}
	return UserToGetUserResponse(user), err
}

func (s *UserService) DeleteUserGrpc(ctx context.Context, request *g.DeleteUserRequest) (*g.DeleteUserResponse, error) {
	info := "grpc service for deleting a user"

	if request == nil {
		s.logger.Warn(info, zap.String("error", errors.StatusInvalidArgumentError), zap.Any("request", request))
		return nil, status.Error(codes.InvalidArgument, errors.StatusInvalidArgumentError)
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultServiceTimeout)
	defer cancel()

	chErr := make(chan error)
	ch := make(chan *struct{})

	go func() {
		defer close(chErr)

		id, err := grpcUuidToUUID(request)
		if err != nil {
			chErr <- err
			return
		}

		err = s.repo.DeleteUser(ctx, id)
		if err != nil {
			chErr <- err
			return
		}
	}()

	_, err := httptools.SelectChannels[struct{}](ctx, ch, chErr, info, request, s.logger)
	return &g.DeleteUserResponse{}, err
}
