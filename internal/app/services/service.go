package services

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Repository defines the interface for interacting with the storage backend.
//
//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Repository interface {
	StoreNewUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, login string, hashedPassword []byte) error
}

type GmartServices struct {
	repository     Repository
	accrualAddress string
	dbPool         *pgxpool.Pool //opened in main func dbPool pool connections
}

func NewGmartServices(repository Repository, accrualAddress string, dbPool *pgxpool.Pool) *GmartServices {
	return &GmartServices{
		repository:     repository,
		accrualAddress: accrualAddress,
		dbPool:         dbPool,
	}
}
