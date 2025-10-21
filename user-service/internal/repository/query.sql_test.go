package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*Queries, func()) {
	ctx := context.Background()

	// Запускаем PostgreSQL в Docker
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	// 2. Получаем строку подключения
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// 3. Подключаемся через pgx
	conn, err := pgx.Connect(ctx, connStr)
	require.NoError(t, err)

	// 4. Создаём схему (вместо миграций — для простоты тестов)
	_, err = conn.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMPTZ
		);
		CREATE UNIQUE INDEX idx_users_email_not_deleted ON users (email) WHERE deleted_at IS NULL;
	`)
	require.NoError(t, err)

	// 5. Создаём querier
	queries := New(conn) // ✅ conn реализует DBTX

	// 6. Cleanup
	cleanup := func() {
		conn.Close(ctx)
		pgContainer.Terminate(ctx)
	}

	return queries, cleanup
}

func TestGetUserByEmail(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	email := "findme@example.com"
	hash := "$2a$10$validhash"

	// Создаём пользователя
	created, err := queries.CreateUser(context.Background(), CreateUserParams{
		Email:        email,
		PasswordHash: hash,
	})
	require.NoError(t, err)

	// Ищем по email
	found, err := queries.GetUserByEmail(context.Background(), email)
	require.NoError(t, err)
	require.Equal(t, created.ID, found.ID)
	require.Equal(t, email, found.Email)
	require.False(t, found.DeletedAt.Valid)
}

func TestDeleteUser(t *testing.T) {
	queries, cleanup := setupTestDB(t)
	defer cleanup()

	email := "todelete@example.com"
	hash := "$2a$10$validhash"

	// Создаём пользователя
	created, err := queries.CreateUser(context.Background(), CreateUserParams{
		Email:        email,
		PasswordHash: hash,
	})
	require.NoError(t, err)

	// Удаляем (soft delete)
	err = queries.DeleteUser(context.Background(), created.ID)
	require.NoError(t, err)

	// Пытаемся найти — должно вернуться sql.ErrNoRows
	_, err = queries.GetUserByEmail(context.Background(), email)
	require.Error(t, err, "expected user not found after soft delete")
	require.Equal(t, pgx.ErrNoRows, err)
}
