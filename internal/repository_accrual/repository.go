package repositoryaccrual

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	m "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/models_accrual"
)

type Repository interface {
	CreateOrder(req *m.CreateOrderRequest) error
	IsExistOrder(orderID string) (bool, error)
	GetOrdersForProcessing() ([]m.Order, error)
	ChangeStatus(orderID string, status m.OrderStatus) error
	GetMatchingReward(good m.Goods) (m.Reward, error)
	AddAccrual(orderID string, accrual float64) error
}

type PgStorage struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("db not init")
	}
	return &PgStorage{
		db: db,
	}, nil
}

func (r *PgStorage) IsExistOrder(orderID string) (bool, error) {
	var exist bool
	err := r.db.QueryRow(queryExistOrderID, orderID).Scan(&exist)
	if err != nil {
		return exist, fmt.Errorf("repo.IsExistOrder() %w", err)
	}

	return exist, nil
}

func (r *PgStorage) CreateOrder(req *m.CreateOrderRequest) error {
	tx, err := r.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = r.db.Exec(querySetOrder, *req.Order, m.OrderStatusNew)
	if err != nil {
		return fmt.Errorf("repo.CreateOrder() %w", err)
	}

	for _, good := range req.Goods {
		_, err := r.db.Exec(querySetOrdeData, *req.Order, good.Description, good.Price)
		if err != nil {
			return fmt.Errorf("repo.CreateOrder() %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PgStorage) GetOrdersForProcessing() ([]m.Order, error) {
	var orders []m.Order

	rows, err := r.db.Query(queryGetOrderForProcessing, m.OrderStatusProcessing)
	if err != nil {
		return nil, fmt.Errorf("repo.GetOrdersForProcessing() %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := m.Order{}
		var accrual sql.NullFloat64
		if err = rows.Scan(&model.ID, &model.Status, &accrual); err != nil {
			return nil, err
		}
		model.Accrual = accrual.Float64
		orders = append(orders, model)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orders, nil
}

func (r *PgStorage) ChangeStatus(orderID string, status m.OrderStatus) error {

	_, err := r.db.Exec(
		queryChangeStatusOrder,
		status,
		orderID,
	)

	return err
}

func (r *PgStorage) GetMatchingReward(good m.Goods) (m.Reward, error) {

	var reward m.Reward
	err := r.db.QueryRow(
		queryGetReward,
		good.Description,
	).Scan(&reward.Reward, &reward.RewardType)

	if errors.Is(err, sql.ErrNoRows) {
		return reward, nil
	}

	return reward, err
}

func (r *PgStorage) AddAccrual(orderID string, accrual float64) error {

	_, err := r.db.Exec(
		queryAddAccrual,
		accrual,
		m.OrderStatusProcessed,
		orderID,
	)

	return err
}