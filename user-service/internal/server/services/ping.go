package services

import (
	"context"
	g "github.com/task_platform/tools/user-service/grpc"

	errors "github.com/task_platform/tools/helpers"
	"go.uber.org/zap"
)

type HealthService struct {
	g.UnimplementedUserServiceHealthGrpcServer
	logger *zap.Logger
}

func NewHealthService(logger *zap.Logger) *HealthService {
	return &HealthService{
		UnimplementedUserServiceHealthGrpcServer: g.UnimplementedUserServiceHealthGrpcServer{},
		logger:                                   logger,
	}
}

func (s *HealthService) PingGrpc(ctx context.Context, req *g.PingRequest) (*g.PingResponse, error) {
	info := "grpc orders health service ping"

	ctx, cancel := context.WithTimeout(ctx, DefaultServiceTimeout)
	defer cancel()

	ch := make(chan bool)

	go func() {
		defer close(ch)

		if req == nil {
			s.logger.Warn(info, zap.String("error", errors.StatusInvalidArgumentError), zap.Any("req", req))
			ch <- false
		}

		ch <- true
	}()

	out := new(g.PingResponse)

	select {
	case <-ctx.Done():
		out.Data = false
		return out, nil
	case res := <-ch:
		out.Data = res
		return out, nil
	}
}
