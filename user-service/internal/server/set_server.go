package server

import (
	"context"
	"errors"
	"github.com/atsybenko962/task-platform/user-service/internal/repository"
	"github.com/atsybenko962/task-platform/user-service/internal/server/services"
	"github.com/task_platform/tools/configcore"
	"github.com/task_platform/tools/grpctools"
	g "github.com/task_platform/tools/user-service/grpc"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Serve(ctx context.Context, config configcore.ServerConfig) error
}

func SetServer(db repository.DBTX, cfg configcore.ServerConfig, logger *zap.Logger) (Server, error) {
	srvGRPC := grpc.NewServer()

	userSrv := services.NewUserService(db, logger)
	healthSrv := services.NewHealthService(logger)

	g.RegisterUserServiceServer(srvGRPC, userSrv)
	g.RegisterUserServiceHealthGrpcServer(srvGRPC, healthSrv)

	serverGRPC := grpctools.NewServerGRPC(srvGRPC, logger)

	logger.Info("user grpc server started", zap.String("host", cfg.ServerHost),
		zap.String("port", cfg.ServerPort), zap.String("protocol", cfg.ServerType))

	return serverGRPC, nil
}

func Run(cfg configcore.ServerConfig, server Server, logger *zap.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())

	g, gCtx := errgroup.WithContext(ctx)

	shutdownChan := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		sigInt := <-sig
		logger.Info("signal interrupt received", zap.Stringer("os_signal", sigInt))
		close(sig)

		shutdownChan <- struct{}{}
		cancel()

		return nil
	})

	g.Go(func() error {
		err := server.Serve(gCtx, cfg)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("app: user server error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Error("wait group: user server error", zap.Error(err))
		return err
	}
	<-shutdownChan
	close(shutdownChan)

	return nil

}
