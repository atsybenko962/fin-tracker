package httptools

import (
	"context"
	errors "github.com/task_platform/tools/helpers"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

func SelectChannels[T any](ctx context.Context, ch chan *T, chErr chan error, info string, in any, logger *zap.Logger) (*T, error) {
	select {
	case <-ctx.Done():
		logger.Warn(info, zap.String("error", "time is out"), zap.Any("in", in))
		return nil, status.Error(codes.DeadlineExceeded, errors.StatusDeadlineExceeded)

	case err := <-chErr:
		if strings.Contains(err.Error(), errors.StatusItemNotFound) {
			logger.Error(info, zap.String("error", "not found"), zap.Error(err), zap.Any("in", in))
			return nil, status.Error(codes.NotFound, err.Error())
		}

		logger.Error(info, zap.String("error", "internal"), zap.Error(err), zap.Any("in", in))
		return nil, status.Error(codes.Internal, err.Error())

	case item := <-ch:
		logger.Debug(info, zap.String("status", "success"))
		return item, nil
	}
}
