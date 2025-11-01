package main

import (
	"context"
	"github.com/atsybenko962/task-platform/user-service/internal/config"
	logs "github.com/atsybenko962/task-platform/user-service/internal/logger"
	"github.com/atsybenko962/task-platform/user-service/internal/server"
	"github.com/task_platform/tools/configcore"
	"github.com/task_platform/tools/dbtools"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	var cfg config.Config
	if err := configcore.Load(cfg, ""); err != nil {
		log.Panic("Can't load config file", err)
	}

	logger := logs.NewLogger(cfg, os.Stdout)

	db, err := dbtools.NewClient(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Panic("user db connection error", zap.Error(err))
	}

	serverGRPC, err := server.SetServer(db, cfg.UserServerConfig, logger)
	if err != nil {
		logger.Panic("set user grpc server error", zap.Error(err))
	}

	err = server.Run(cfg.UserServerConfig, serverGRPC, logger)
	if err != nil {
		logger.Panic("run user grpc server error", zap.Error(err))
	}

}
