package helpers

import (
	"database/sql"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
)

func GetUserBalance(row *sql.Row) (user.Balance, error) {
	var userBalance user.Balance
	if row.Err() != nil {
		return user.Balance{}, row.Err()
	}
	if err := row.Scan(&userBalance.UserID, &userBalance.Current, &userBalance.Withdrawn); err != nil {
		return user.Balance{}, err
	}
	return userBalance, nil
}

func GetOrders(rows *sql.Rows) ([]user.Order, error) {
	var orders []user.Order
	for rows.Next() {
		var order user.Order
		var accrual sql.NullFloat64
		err := rows.Scan(&order.Number, &accrual, &order.UserID, &order.Status, &order.UploadedAt, &order.ProcessedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		orders = append(orders, order)
	}
	return orders, nil
}
