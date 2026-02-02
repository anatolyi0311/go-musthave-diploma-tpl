package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/helpers"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/storage"
)

type accrual struct {
	client          *http.Client
	accrualAddress  string
	dbStorage       storage.GofferStorage
	intervalCall    int
	activeJobsMutex sync.RWMutex
	activeJobs      map[string]context.CancelFunc
	jobQueue        chan jobCommand
	stopCh          chan struct{}
	pauseUntil      time.Time
	pauseUntilMu    sync.RWMutex
}

type accrualOrder struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type jobCommand struct {
	OrderNumber string
	UserID      int
	Done        chan struct{} // Сигнал завершения
}

func NewAccrual(host string, dbStorage storage.GofferStorage) *accrual {
	return &accrual{
		client:         &http.Client{Timeout: 10 * time.Second},
		accrualAddress: host + "/api/orders/",
		dbStorage:      dbStorage,
		intervalCall:   1,
		activeJobs:     make(map[string]context.CancelFunc),
		jobQueue:       make(chan jobCommand, 100),
		stopCh:         make(chan struct{}),
		pauseUntil:     time.Time{},
	}

}

func (a *accrual) getOrderAccrual(order string) (*accrualOrder, error) {

	var accrualOrder *accrualOrder

	response, err := a.client.Get(a.accrualAddress + order)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//429 — превышено количество запросов к сервису.
	if response.StatusCode == http.StatusTooManyRequests {
		a.pauseUntilMu.Lock()
		a.pauseUntil = time.Now().Add(60 * time.Second)
		a.pauseUntilMu.Unlock()

		return nil, fmt.Errorf("received 429, pausing all requests for 60 seconds")
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&accrualOrder); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return accrualOrder, nil
}

func (a *accrual) getOpenOrders() {

	query := `SELECT * FROM orders WHERE status IN ('NEW', 'PROCESSING')`
	rows, err := a.dbStorage.DBStorage.GetRows(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	orders, _ := helpers.GetOrders(rows)
	for _, order := range orders {
		a.activeJobsMutex.RLock()
		running := a.activeJobs[order.Number] != nil
		a.activeJobsMutex.RUnlock()

		if running {
			continue
		}

		done, err := a.addOrder(order.Number, order.UserID)
		if err != nil {
			fmt.Printf("Failed to start job for order %s: %v\n", order.Number, err)
		} else {
			fmt.Printf("Started job for order %s from poller\n", order.Number)
			_ = done
		}
	}
}

func (a *accrual) addOrder(orderNumber string, userID int) (<-chan struct{}, error) {
	a.activeJobsMutex.Lock()
	defer a.activeJobsMutex.Unlock()

	if _, exists := a.activeJobs[orderNumber]; exists {
		return nil, fmt.Errorf("job for order %s is already running", orderNumber)
	}

	done := make(chan struct{})
	cmd := jobCommand{
		OrderNumber: orderNumber,
		UserID:      userID,
		Done:        done,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	a.activeJobs[orderNumber] = cancel

	go func() {
		select {
		case a.jobQueue <- cmd:
			fmt.Printf("Job enqueued for order %s\n", orderNumber)
			go a.runJob(cmd)
		case <-ctx.Done():
			a.cleanupJob(orderNumber)
			close(done)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Printf("Job for order %s timed out after 1 hour\n", orderNumber)
			}
			close(done)
		case <-done:
		}
	}()

	return done, nil
}

func (a *accrual) runJob(cmd jobCommand) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.activeJobsMutex.Lock()
	if existingCancel, exists := a.activeJobs[cmd.OrderNumber]; exists {
		defer existingCancel()
	}
	a.activeJobsMutex.Unlock()

	defer func() {
		a.cleanupJob(cmd.OrderNumber)
	}()

	delay := time.Second
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			a.pauseUntilMu.RLock()
			pauseUntil := a.pauseUntil
			a.pauseUntilMu.RUnlock()

			if time.Now().Before(pauseUntil) {
				fmt.Printf("Skipping check for %s: service is paused due to rate limit\n", cmd.OrderNumber)
				continue
			}

			accrualData, err := a.getOrderAccrual(cmd.OrderNumber)
			if err != nil {
				fmt.Printf("Error fetching accrual for %s: %v\n", cmd.OrderNumber, err)
				delay += time.Second
				ticker.Reset(delay)
				continue
			}

			if accrualData == nil {
				// 204 No Content — заказ не зарегистрирован в системе расчёта
				delay += time.Second
				ticker.Reset(delay)
				continue
			}

			// Обновляем статус
			var query string
			if accrualData.Status == "PROCESSED" || accrualData.Status == "INVALID" {
				query = `UPDATE orders SET status = $1, accrual = $2, processed_at = $4 WHERE number = $3`
				err := a.dbStorage.DBStorage.Insert(query,
					accrualData.Status, accrualData.Accrual, cmd.OrderNumber, time.Now())
				if err != nil {
					fmt.Printf("Failed to update order %s: %v\n", cmd.OrderNumber, err)
				}

				if accrualData.Status == "PROCESSED" {
					err := a.updateOrderOperation(accrualData, cmd.UserID)
					if err != nil {
						fmt.Printf("Fail update order %s in balance: %v\n", cmd.OrderNumber, err)
					}
				}

				fmt.Printf("Order %s reached final status: %s\n", cmd.OrderNumber, accrualData.Status)
				return // Завершаем job
			} else {
				fmt.Printf("Order %s is %s\n", cmd.OrderNumber, accrualData.Status)
				query = `UPDATE orders SET status = $1, processed_at = $3 WHERE number = $2`
				_ = a.dbStorage.DBStorage.Insert(query,
					accrualData.Status, cmd.OrderNumber, time.Now())
				delay += time.Second
				ticker.Reset(delay)
				continue
			}

		case <-ctx.Done():
			fmt.Printf("Job for order %s stopped: %v\n", cmd.OrderNumber, ctx.Err())
			return
		}
	}
}

func (a *accrual) updateOrderOperation(order *accrualOrder, userID int) error {

	var query string

	//логгируем операцию
	query = `INSERT INTO balance_operations (user_id, amount, "type", order_number, processed_at) VALUES ($1, $2, $3, $4, $5)`
	if err := a.dbStorage.DBStorage.Insert(query, userID, order.Accrual, "accrual", order.Order, time.Now()); err != nil {
		return err
	}

	//сохраняем баланс
	query = `
		UPDATE user_balance 
		SET current = current + $2, 
			updated_at = $3 
		WHERE user_id = $1
	`
	err := a.dbStorage.DBStorage.Insert(query, userID, order.Accrual, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (a *accrual) cleanupJob(orderNumber string) {
	a.activeJobsMutex.Lock()
	defer a.activeJobsMutex.Unlock()

	if cancel, exists := a.activeJobs[orderNumber]; exists {
		cancel()
		delete(a.activeJobs, orderNumber)
	}
}

func (a *accrual) Start() {
	//TODO: подумать как создавать JOB при получении заказа, а не пушить БД каждые 5 секунд
	ticker := time.NewTicker(5 * time.Second) //проверяем БД каждый 5 секунд
	defer ticker.Stop()

	for {
		select {
		case <-a.stopCh:
			return
		case <-ticker.C:
			go a.getOpenOrders()
		}
	}
}

func (a *accrual) Stop() {
	close(a.stopCh)
}
