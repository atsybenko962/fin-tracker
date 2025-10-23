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
	panic("Implement me")
}

func (s *UserService) DeleteUserGrpc(ctx context.Context, request *g.DeleteUserRequest) (*g.DeleteUserResponse, error) {
	panic("Implement me")
}
