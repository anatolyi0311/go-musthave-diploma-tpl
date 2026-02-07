package server_accrual

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	config "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/config_accrual"
	m "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/models_accrual"
	uc "github.com/anatolyi0311/go-musthave-diploma-tpl/internal/usecase_accrual"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Server struct {
	cfg  *config.Config
	core *chi.Mux
	uc   uc.UseCase
}

func NewServer(cfg *config.Config, db *sql.DB) (*Server, error) {
	uc, err := uc.NewService(db)
	if err != nil {
		return nil, err
	}

	server := &Server{
		core: chi.NewRouter(),
		uc:   uc,
		cfg:  cfg,
	}
	server.route()

	return server, nil
}

func (s *Server) route() {

	s.core.Post("/api/orders", s.CreateOrder)
}

func (s *Server) RunEngine(ctx context.Context) {
	s.uc.RunEngine(ctx)
}

func (s *Server) Run() error {
	logrus.Infof("server started with params: host - %s", s.cfg.Opts.ServerHost)
	return http.ListenAndServe(s.cfg.Opts.ServerHost, s.core)
}

func (s *Server) CreateOrder(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method must be POST", http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(res, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "cannot read body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var request *m.CreateOrderRequest
	if err = json.Unmarshal(body, &request); err != nil {
		logrus.Errorln(err)
		http.Error(res, "cannot unmarshal body", http.StatusBadRequest)
		return
	}

	status := http.StatusAccepted
	err = s.uc.CreateOrder(request)
	if err != nil {
		if errors.Is(err, m.OrderAlreadyExists) {
			status = http.StatusConflict
		} else {
			status = http.StatusBadRequest
		}
	}
	res.WriteHeader(status)
}
