package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework3/cmd/app"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "9999"
const defaultHost = "0.0.0.0"
const clientsDb = "postgres://app:pass@clientsdb:5432/db"
const suggestionDb = "mongodb://app:pass@suggestiondb:27017/db"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	if err := execute(net.JoinHostPort(host, port), clientsDb, suggestionDb); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string, cliDb string, sugDb string) error {

	mux := chi.NewMux()

	ctxCLiDb := context.Background()
	poolCli, err := pgxpool.Connect(ctxCLiDb, clientsDb)
	if err != nil {
		return err
	}
	defer poolCli.Close()

	ctxSugDb := context.Background()
	poolSug, err := pgxpool.Connect(ctxSugDb, sugDb)
	if err != nil {
		return err
	}
	defer poolSug.Close()

	application := app.NewServer(mux, poolCli, ctxCLiDb, poolSug, ctxSugDb)
	err = application.Init()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}
