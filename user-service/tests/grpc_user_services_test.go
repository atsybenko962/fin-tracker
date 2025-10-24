package tests

import (
	"context"
	"github.com/atsybenko962/task-platform/user-service/internal/server/services"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	g "github.com/task_platform/tools/user-service/grpc"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

func setupGrpcClient(t *testing.T) (*grpc.ClientConn, func()) {
	ctx := context.Background()

	//Запуск тестовой БД в Docker
	pgConteiner, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	//Получение строки подключения
	connStr, err := pgConteiner.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	//Подключение к БД
	db, err := pgx.Connect(ctx, connStr)
	require.NoError(t, err)

	//Создание таблицы
	_, err = db.Exec(ctx, `
		CREATE TABLE users (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    		email TEXT NOT NULL UNIQUE,
    		password_hash TEXT NOT NULL,
    		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    		deleted_at TIMESTAMP
		);
`)
	require.NoError(t, err)

	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	logger, _ := zap.NewDevelopment()
	grpcServer := grpc.NewServer()

	nus := services.NewUserService(db, logger)

	g.RegisterUserServiceServer(grpcServer, nus)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	closeconn := func() {
		conn.Close()
		grpcServer.Stop()
		db.Close(ctx)
		pgConteiner.Terminate(ctx)
	}

	return conn, closeconn
}
