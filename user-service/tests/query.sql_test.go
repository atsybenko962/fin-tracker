package tests

import (
	"context"
	rep "github.com/atsybenko962/task-platform/user-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*rep.Queries, func()) {
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
	conn, err := pgx.Connect(ctx, connStr)
	require.NoError(t, err)

	//Создание таблицы
	_, err = conn.Exec(ctx, `
		CREATE TABLE users (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    		email TEXT NOT NULL UNIQUE,
    		password_hash TEXT NOT NULL,
    		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    		deleted_at TIMESTAMP
		);
`)
	require.NoError(t, err)

	//
	queries := rep.New(conn)

	// Закрытие соединения

	closeconn := func() {
		conn.Close(ctx)
		pgConteiner.Terminate(ctx)
	}

	return queries, closeconn
}

func TestUserReposytory(t *testing.T) {
	queries, closeconn := setupTestDB(t)
	defer closeconn()

	email := "test@test.ru"
	hash := "test1234hash1234"

	testID := 0
	t.Logf("\tTest %d:\t Testing user creation", testID)

	user, err := queries.CreateUser(context.Background(), rep.CreateUserParams{Email: email, PasswordHash: hash})
	require.NoError(t, err)
	require.Equal(t, email, user.Email)
	require.Equal(t, hash, user.PasswordHash)

	testID++
	t.Logf("\tTest %d:\t Testing the user's email search", testID)

	{
		foundUser, err := queries.GetUserByEmail(context.Background(), email)
		require.NoError(t, err)

		require.Equal(t, user.ID, foundUser.ID)
		require.Equal(t, user.Email, foundUser.Email)
		require.False(t, foundUser.DeletedAt.Valid)
	}

	testID++
	t.Logf("\tTest %d:\t Testing a user's soft deletion", testID)

	{
		err = queries.DeleteUser(context.Background(), user.ID)
		require.NoError(t, err)

		_, err = queries.GetUserByEmail(context.Background(), email)
		require.Equal(t, pgx.ErrNoRows, err)
	}
}
