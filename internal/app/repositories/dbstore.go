package repositories

// TODO разделить репозиторий на отдельные репозитории

import (
	"context"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"

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

// StoreNewUser сохраняет нового пользователя (заранее сгенерированный UUID, логин и хешированный пароль)
func (d *InDBRepo) StoreNewUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, login string, hashedPassword []byte) error {

	const sqlQuery = `INSERT INTO users (uuid,login, hashed_password) VALUES ($1, $2,$3) ON CONFLICT (login) DO NOTHING`
	_, err := tx.Exec(ctx, sqlQuery, userID, login, hashedPassword)
	if err != nil {
		logrus.Error("new user don't save in database ", err)
		return err
	}
	return nil
}

// StoreNewUserBalance создает поле с нулевым балансом для пользователя в таблице balance
func (d *InDBRepo) StoreNewUserBalance(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	const sqlQuery = `INSERT INTO balance (uuid) VALUES ($1)`
	_, err := tx.Exec(ctx, sqlQuery, userID)
	if err != nil {
		logrus.Error("user balance (0.00) don't save in database ", err)
		return err
	}
	return nil
}
