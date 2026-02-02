package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type ContextKey string

const UserIDKey ContextKey = "userID"

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"-"`
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	Number      string     `json:"number"`
	Status      string     `json:"status"`
	Accrual     *float64   `json:"accrual,omitempty"`
	UploadedAt  time.Time  `json:"uploaded_at"`
	ProcessedAt *time.Time `json:"-"`
	UserID      int        `json:"-"`
}

type Balance struct {
	UserID    int     `json:"-"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

type Operation struct {
	ID          int       `json:"-"`
	UserID      int       `json:"-"`
	Type        string    `json:"-"`
	Amount      float64   `json:"sum"`
	Order       string    `json:"order"`
	ProcessedAt time.Time `json:"processed_at"`
}

type RequestOrder struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}