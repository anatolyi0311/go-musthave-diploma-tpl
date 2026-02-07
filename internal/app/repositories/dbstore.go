package repositories

// TODO разделить репозиторий на отдельные репозитории

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type InDBRepo struct {
	dbPool *pgxpool.Pool //opened in main func dbPool pool connections
}

func NewURLInDBRepo(dbPool *pgxpool.Pool) *InDBRepo {
	storage := &InDBRepo{
		dbPool: dbPool,
	}
	if err := storage.CreateBDTables(); err != nil {
		logrus.Error(err)
	}
	return storage
}

func (d *InDBRepo) CreateBDTables() error {
	return nil
}
