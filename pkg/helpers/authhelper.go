package helpers

import (
	"net/http"
	"time"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/user"
	"github.com/golang-jwt/jwt/v4"
)

func GetAuthCookie(userID int, login string, jwtKey []byte) (*http.Cookie, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &user.Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{ // ← изменилось имя поля
			ExpiresAt: jwt.NewNumericDate(expirationTime), // ← новый способ
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
		Path:    "/",
	}

	return &cookie, nil
}