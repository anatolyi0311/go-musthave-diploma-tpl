package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/crypto"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
)

func (d *DataBase) RegisterUserWithBalance(req user.UserRequest, isHashed bool) (int, error) {

	transaction, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	defer transaction.Rollback()

	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password) // хеширование пароля
	}

	var userID int
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	err = transaction.QueryRow(query, req.Login, pass).Scan(&userID)
	if err != nil {
		return 0, err
	}

	t := time.Now()
	_, err = transaction.Exec(`INSERT INTO user_balance (user_id, current, withdrawn, updated_at, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5)`, userID, 0, 0, t, t)
	if err != nil {
		return 0, err
	}

	// Фиксируем транзакцию
	if err = transaction.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}

func (d *DataBase) UserIsRegistred(login string) (bool, error) {
	query := `SELECT COUNT(id) FROM users WHERE login = $1`
	result, err := d.CountRows(query, login)
	return result > 0, err
}

func (d *DataBase) GetUserBalance(userID int) *sql.Row {
	query := `SELECT user_id, current, withdrawn FROM user_balance WHERE user_id = $1`
	return d.InsertWithReturning(query, userID)
}

func (d *DataBase) GetOrders(r context.Context, userID int) (*sql.Rows, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1 ORDER BY uploaded_at DESC`
	return d.GetRows(r, query, userID)
}

func (d *DataBase) GetWithdrawals(r context.Context, userID int) (*sql.Rows, error) {
	query := `SELECT * FROM balance_operations WHERE user_id = $1 AND type = 'withdraw'  ORDER BY processed_at DESC`
	return d.GetRows(r, query, userID)
}

func (d *DataBase) GetUserID(req user.UserRequest, isHashed bool) (int, error) {
	query := `SELECT id FROM users WHERE login = $1 AND password = $2`
	pass := req.Password
	if !isHashed {
		pass = crypto.HashString(req.Password)
	}
	return d.CountRows(query, req.Login, pass)
}

func (d *DataBase) GetUserIDWithOrder(orderNumber string) (int, error) {
	query := `SELECT user_id FROM orders WHERE number = $1`
	return d.CountRows(query, orderNumber)
}

func (d *DataBase) AddNewOrder(orderNumber string, userID int, statusOrder string) error {
	query := `INSERT INTO orders (number, user_id, status, uploaded_at) VALUES ($1, $2, $3, $4)`
	t := time.Now()
	return d.Insert(query, orderNumber, userID, statusOrder, t)
}

func (d *DataBase) UpdateBalance(userID int, req user.RequestOrder, userBalance user.Balance) error {
	transaction, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	t := time.Now()

	query := `
		UPDATE user_balance
		SET current = $2,
			withdrawn = $3,
			updated_at = $4
		WHERE user_id = $1
	`
	_, err = transaction.Exec(query, userID, userBalance.Current, userBalance.Withdrawn, t)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO balance_operations (user_id, type, amount, order_number, processed_at)
		VALUES ($1, 'withdraw', $2, $3, $4)
	`
	_, err = transaction.Exec(query, userID, req.Sum, req.Order, t)
	if err != nil {
		return err
	}

	return transaction.Commit()
}
