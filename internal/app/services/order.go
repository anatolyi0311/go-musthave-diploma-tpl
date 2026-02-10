package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks
type Order interface {
	StoreUserOrder(ctx context.Context, tx pgx.Tx, orderNumber, orderStatus string, userID uuid.UUID, bonus decimal.Decimal) error
}
