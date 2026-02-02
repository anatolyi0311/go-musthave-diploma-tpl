package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/config"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/handler"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/middleware"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/accrual"
	"github.com/go-chi/chi"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config") // .Err(err).Msg("Failed to load config")
	}

	router := chi.NewRouter()
	h := handler.NewHandler(cfg)

	router.Use(middleware.LoggingMiddleware)

	router.Post("/api/user/register", h.RegisterUserHandler)
	router.Post("/api/user/login", h.LoginHandler)

	router.Group(func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(cfg.SecretKey))

		// router.Post("/api/user/orders", h.UploadOrderHandler)
		// router.Get("/api/user/orders", h.GetOrdersHandler)

		// router.Get("/api/user/balance", h.GetBalanceHandler)

		// router.Post("/api/user/balance/withdraw", h.WithdrawHandler)
		// router.Get("/api/user/withdrawals", h.GetWithdrawalsHandler)
	})

	accrual := accrual.NewAccrual(cfg.AccrualAddress, h.Storage)
	go accrual.Start()
	log.Printf("Accrual service available on: %s...", cfg.AccrualAddress)

	srv := &http.Server{
		Addr:    cfg.Host,
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on: %s...", cfg.Host)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error") //.Err(err).Msg("Server error")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Graceful shutdown initiated...") //.Msg("Graceful shutdown initiated...")

	// Контекст для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Останавливаем accrual-сервис
	accrual.Stop()

	// Останавливаем HTTP-сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Forced server shutdown") //.Err(err).Msg("Forced server shutdown")
	} else {
		log.Println("Server stopped gracefully") //.Info().Msg("Server stopped gracefully")
	}
}
