package main

import (
	"context"
	"errors"
	"os"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/app/config"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/app/handlers"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/app/repositories"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/app/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"net/http"
)

func init() {
	// Установка флага для сериализации decimal в JSON без кавычек
	decimal.MarshalJSONWithoutQuotes = true
}

func main() {
	var (
		dbPool               *pgxpool.Pool
		err                  error
		cfg                  *config.ENVConfig
		GophermartRepository services.Repository
	)

	cfg = config.NewConfig()
	confPool, err := pgxpool.ParseConfig(cfg.EnvDataBase)
	if err != nil {
		logrus.Errorf("error parsing config: %v", err)
	}
	confPool.MaxConns = 50
	confPool.MinConns = 10
	dbPool, err = pgxpool.NewWithConfig(context.Background(), confPool)
	if err != nil {
		logrus.Error("Don't connect to dbPool: ", err)
		os.Exit(1)
	}

	defer dbPool.Close()
	GophermartRepository = repositories.NewURLInDBRepo(dbPool)

	GophermartService := services.NewGmartServices(GophermartRepository, cfg.EnvAccrualSystemAddress, dbPool)
	GophermartHandler := handlers.NewHandlers(GophermartService, dbPool)

	router := gin.Default()
	publicRoutes := router.Group("/api/user")
	publicRoutes.Use(GophermartHandler.MiddlewareLogging())

	server := &http.Server{Addr: cfg.EnvServAdr, Handler: router}
	logrus.Info("Starting server on: ", cfg.EnvServAdr)

	go func() {
		if err = server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logrus.Error(err)
		}
	}()
}
