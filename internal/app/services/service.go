package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface{}

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
