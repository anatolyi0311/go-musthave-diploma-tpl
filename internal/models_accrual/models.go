package models_accrual

type CreateOrderRequest struct {
	Order *string `json:"order"`
	Goods []Goods `json:"goods"`
}

type Goods struct {
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
}

type Order struct {
	Items     []Goods
	ID        string
	Status    OrderStatus
	Accrual   float64
}
