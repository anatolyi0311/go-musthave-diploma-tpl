package usecase_accual

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	engine "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/accrual"
	m "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/models_accrual"
	repo "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/repository_accrual"
)

type UseCase interface {
	RunEngine(ctx context.Context)
	CreateOrder(req *m.CreateOrderRequest) error
}

type Usecase struct {
	repo   repo.Repository
	engine *engine.AccrualEngine
}

func NewService(db *sql.DB) (UseCase, error) {
	repo, err := repo.NewRepo(db)
	if err != nil {
		return nil, fmt.Errorf("fail to init repo")
	}

	engine := engine.New(repo)
	return &Usecase{
		repo:   repo,
		engine: engine,
	}, nil
}

func (u *Usecase) RunEngine(ctx context.Context) {
	u.engine.StartProcessing(ctx)
}

func (u *Usecase) CreateOrder(req *m.CreateOrderRequest) error {
	if req.Order == nil || *req.Order == "" {
		return errors.New("uc.CreateOrder() incorrect orderID")
	}

	if len(req.Goods) < 1 {
		return errors.New("uc.CreateOrder() len(req.Goods) < 1")
	}

	exist, err := u.repo.IsExistOrder(*req.Order)
	if err != nil {
		return err
	}
	if exist {
		return m.OrderAlreadyExists
	}

	for _, good := range req.Goods {
		if !validate(good) {
			return errors.New("uc.CreateOrder() incorrect good data")
		}
	}

	err = u.CreateOrder(req)
	if err == nil {
		order := m.Order{
			Items: req.Goods,
			ID:    *req.Order,
		}
		go u.engine.RegisterOrderForProcessing(order)
	}
	return err
}

func validate(good m.Goods) bool {
	if good.Description == nil || good.Price == nil {
		return false
	}
	if *good.Description == "" || *good.Price < 0 {
		return false
	}
	return true
}
