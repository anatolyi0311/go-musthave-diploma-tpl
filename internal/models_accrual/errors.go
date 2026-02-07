package models_accrual

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	OrderAlreadyExists Err = "order already exists"
)