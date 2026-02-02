package handler

import (
	"github.com/anatolyi0311/go-musthave-diploma-tpl/internal/config"
	"github.com/anatolyi0311/go-musthave-diploma-tpl/pkg/storage"
)

type Handler struct {
	Storage storage.GofferStorage
	JwtKey  []byte
}

func NewHandler(cfg *config.ServerConfig) *Handler {
	storage := storage.NewGofferStorage(cfg.DataBaseDsn)
	return &Handler{
		Storage: *storage,
		JwtKey:  []byte(cfg.SecretKey),
	}
}
