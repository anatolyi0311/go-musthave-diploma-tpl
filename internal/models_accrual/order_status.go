package models_accrual

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "REGISTERED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusError      OrderStatus = "PROCESSING"
)
