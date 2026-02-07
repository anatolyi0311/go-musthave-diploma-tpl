package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"
	"time"
)

//go:generate mockgen -source=handlers.go -destination=mocks/handlers_mock.go -package=mocks
type Service interface {
}

type Handlers struct {
	service Service
	DB      *pgxpool.Pool
}
type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func NewHandlers(service Service, DB *pgxpool.Pool) *Handlers {
	return &Handlers{
		service: service,
		DB:      DB,
	}
}

// MiddlewareLogging provides a logging middleware for Gin.
// It logs details about each request including the URL, method, response status, duration, and size.
// This middleware is useful for monitoring and debugging purposes.
func (h Handlers) MiddlewareLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запуск таймера
		start := time.Now()

		// Обработка запроса
		c.Next()

		// Измерение времени обработки
		duration := time.Since(start)

		// Получение статуса ответа и размера
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Логирование информации о запросе
		logrus.WithFields(logrus.Fields{
			"url":      c.Request.URL.RequestURI(),
			"method":   c.Request.Method,
			"status":   status,
			"duration": duration,
			"size":     size,
		}).Info("Обработан запрос")
	}
}
