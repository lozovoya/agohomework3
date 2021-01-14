package app

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework3/cmd/app/md"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var (
	ErrServer = errors.New("internal server error")
)

type Server struct {
	mux     chi.Router
	poolCli *pgxpool.Pool
	ctxCli  context.Context
	dbSug   *mongo.Database
	ctxSug  context.Context
}

func NewServer(mux chi.Router, poolCli *pgxpool.Pool, ctxCli context.Context, dbsug *mongo.Database, ctxSug context.Context) *Server {
	return &Server{mux: mux, poolCli: poolCli, ctxCli: ctxCli, dbSug: dbsug, ctxSug: ctxSug}
}

func (s *Server) Init() error {

	logMd := middleware.Logger
	identMd := md.IdentMD
	roleCheckerMd := md.IsRole
	authMd := md.AuthMD

	s.mux.With(logMd).Get("/payments", s.Payments)
	s.mux.With(logMd, identMd, roleCheckerMd("USER", s.poolCli, s.ctxCli),
		authMd(s.poolCli, s.ctxCli)).Post("/payments", s.Payments)

	return nil
}

func (s *Server) Payments(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("payments"))
}
