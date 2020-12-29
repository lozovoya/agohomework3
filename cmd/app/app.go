package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

var (
	ErrServer = errors.New("internal server error")
)

type Server struct {
	mux     chi.Router
	poolCli *pgxpool.Pool
	ctxCli  context.Context
	poolSug *pgxpool.Pool
	ctxSug  context.Context
}

func NewServer(mux chi.Router, poolCli *pgxpool.Pool, ctxCli context.Context, poolSug *pgxpool.Pool, ctxSug context.Context) *Server {
	return &Server{mux: mux, poolCli: poolCli, ctxCli: ctxCli, poolSug: poolSug, ctxSug: ctxSug}
}

func (s *Server) Init() error {

	logMd := middleware.Logger

	s.mux.With(logMd).Get("/payments", s.Payments)

	return nil
}

func (s *Server) Payments(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("payments"))
}
