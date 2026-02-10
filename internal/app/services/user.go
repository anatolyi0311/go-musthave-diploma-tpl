package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type User interface {
	StoreNewUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, login string, hashedPassword []byte) error
	GetUserHashPassword(ctx context.Context, login string) ([]byte, error)
	GetUUIDFromUsers(ctx context.Context, login string) (uuid.UUID, error)
}
