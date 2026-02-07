package accrual

import (
	"context"

	m "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/models_accrual"
	repo "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/repository_accrual"
	 		
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type AccrualEngine struct {
	repo              repo.Repository
	ordersToProcessCh chan m.Order
	stopCh            chan struct{}
}

func New(repo repo.Repository) *AccrualEngine {
	return &AccrualEngine{
		repo:              repo,
		ordersToProcessCh: make(chan m.Order),
		stopCh:            make(chan struct{}),
	}
}

func (p *AccrualEngine) RegisterOrderForProcessing(order m.Order) {
	for {
		select {
		case <-p.stopCh:
			logrus.Info("engine stopped")
			return
		case p.ordersToProcessCh <- order:
			logrus.Infof("%s registred for processing", order.ID)
			return
		}
	}
}

func (e *AccrualEngine) StartProcessing(ctx context.Context) {
	logrus.Info("engine running ...")
	newOrders, err := e.repo.GetOrdersForProcessing()
	if err != nil {
		logrus.Warn("cant fetch orders for processing")
	}

	for _, newOrder := range newOrders {
		go e.AddAccrualForOrder(ctx, newOrder)
	}

	for {
		select {
		case <-ctx.Done():
			close(e.stopCh)
			return
		case orderToProcess := <-e.ordersToProcessCh:
			go e.AddAccrualForOrder(ctx, orderToProcess)

		}
	}
}

func (e *AccrualEngine) AddAccrualForOrder(ctx context.Context, order m.Order) {
	err := e.repo.ChangeStatus(order.ID, m.OrderStatusProcessing)
	if err != nil {
		logrus.Error(err, "failed to change status order", order.ID)
		e.markOrderAsFailed(order)
		return
	}

	g, ctx := errgroup.WithContext(ctx)

	results := make([]float64, len(order.Items))

	for i, orderItem := range order.Items {
		i, orderItem := i, orderItem
		g.Go(func() error {
			result, errCalculation := e.calculateAccrualForOrderItem(ctx, orderItem)
			if errCalculation == nil {
				results[i] = result
			}
			return errCalculation
		})
	}

	if err = g.Wait(); err != nil {
		logrus.Error("AddAccrualForOrder(). some of order items not calculated properly")
		e.markOrderAsFailed(order)
		return
	}

	totalAccrual := 0.0
	for _, accrualForItem := range results {
		totalAccrual += accrualForItem
	}

	err = e.repo.AddAccrual(order.ID, totalAccrual)
	if err != nil {
		logrus.Error(err, "failed to change status order", order.ID)
		e.markOrderAsFailed(order)
	}
}

func (e *AccrualEngine) markOrderAsFailed(order m.Order) {
	err := e.repo.ChangeStatus(order.ID, m.OrderStatusError)
	if err != nil {
		logrus.Warn("failed to change status", order.ID)
	}
}

func (e *AccrualEngine) calculateAccrualForOrderItem(_ context.Context, good m.Goods) (float64, error) {
	matchingReward, err := e.repo.GetMatchingReward(good)
	return matchingReward.CalculateReward(*good.Price), err
}