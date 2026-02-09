package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Claims — claims structure that includes standard claims and UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const (
	TokenExp = time.Hour * 3
)

// BuildJWTString creates a token with the HS256 signature algorithm and Claims statements and returns it as a string.
func BuildJWTString(userID uuid.UUID, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.New("Secret key is empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// time create token
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	// create token string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return tokenString, err
}

// GenerateUniqueID генерирует UUID при помощи библиотеки golang.org/x/crypto/bcrypt
func GenerateUniqueID() uuid.UUID {
	return uuid.New()
}

// GetUserID we check the validity of the token and if it is valid, then we get and return the UserID from it
func GetUserID(tokenString, secretKey string) (uuid.UUID, error) {
	if secretKey == "" {
		return uuid.Nil, errors.New("Secret key is empty")
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		logrus.Error(err)
		return uuid.Nil, err
	}
	if !token.Valid {
		err = fmt.Errorf("token is not valid")
		logrus.Error(err)
		return uuid.Nil, err
	}
	logrus.Infof("Token is valid, userID: %v", claims.UserID)
	return claims.UserID, nil
}

// CreateHashPassword хеширует пароль пользователя для сохранения в репозитории
func CreateHashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return hashedPassword, nil
}

// CheckHashPasswordForValid проверяет введенный пользователем пароль на соответствие сохраненному
func CheckHashPasswordForValid(savedHashedPassword []byte, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword(savedHashedPassword, []byte(inputPassword))
	if err != nil {
		logrus.Error(err)
		return false
	} else {
		return true
	}
}
