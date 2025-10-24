package grpctools

import (
	"context"
	"fmt"
	"github.com/task_platform/tools/configcore"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type ServerGRPC struct {
	logger *zap.Logger
	srv    *grpc.Server
}

func NewServerGRPC(srv *grpc.Server, logger *zap.Logger) *ServerGRPC {
	return &ServerGRPC{
		logger: logger,
		srv:    srv,
	}
}

func (s *ServerGRPC) Serve(ctx context.Context, cfg configcore.ServerConfig) error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServerPort))
	if err != nil {
		s.logger.Error("grpc server listen error", zap.Error(err))
		return err
	}

	// запустим Graceful Stop через отдельную горутину
	go func() {
		<-ctx.Done()
		s.logger.Warn("grpc server: shutting down...")
		s.srv.GracefulStop()
	}()

	return s.srv.Serve(l)
}
