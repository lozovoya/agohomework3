package md

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework3/cmd/app/dto"
	"log"
	"net/http"
)

type contextKey struct {
	name string
}

var identifierContextKey = &contextKey{"identifier context"}
var userIdContextKey = &contextKey{"user id"}

func (c *contextKey) String() string {
	return c.name
}

func IdentMD(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var token *dto.TokenDTO
		err := json.NewDecoder(r.Body).Decode(&token)
		if err != nil {
			log.Println(err)
		}

		if token.Token != "" {
			ctx := context.WithValue(r.Context(), identifierContextKey, &token.Token)
			r = r.WithContext(ctx)
			log.Println(token.Token)
			handler.ServeHTTP(w, r)
		}

		//w.WriteHeader(500)
	})
}

func AuthMD(poolCli *pgxpool.Pool, ctxCli context.Context) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, ok := r.Context().Value(identifierContextKey).(*string)
			if !ok {
				log.Println("token error")
				w.WriteHeader(500)
				return
			}

			var user struct {
				login string
			}

			err := poolCli.QueryRow(ctxCli,
				"SELECT login FROM users WHERE token = $1", token).Scan(&user.login)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			if user.login != "" {
				ctx := context.WithValue(r.Context(), userIdContextKey, &user.login)
				r = r.WithContext(ctx)
				log.Println(user.login)
				handler.ServeHTTP(w, r)
			}

			w.WriteHeader(500)
			return
		})
	}
}

func IsRole(role string, poolCli *pgxpool.Pool, ctxCli context.Context) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token, ok := r.Context().Value(identifierContextKey).(*string)
			if !ok {
				log.Println("token error")
				w.WriteHeader(500)
				return
			}

			var user struct {
				role string
			}

			err := poolCli.QueryRow(ctxCli,
				"SELECT role FROM users WHERE token = $1", token).Scan(&user.role)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				return
			}

			if user.role != role {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
