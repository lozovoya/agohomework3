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

func GetToken(r *http.Request) string {
	token, _ := r.Context().Value(identifierContextKey).(*string)
	return *token
}

func GetUserId(r *http.Request) int {
	userid, _ := r.Context().Value(userIdContextKey).(*int)
	return *userid
}

func IdentMD(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var token *dto.TokenDTO
		err := json.NewDecoder(r.Body).Decode(&token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if token.Token == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), identifierContextKey, &token.Token)
		r = r.WithContext(ctx)
		log.Println(token)
		handler.ServeHTTP(w, r)
	})
}

func AuthMD(poolCli *pgxpool.Pool, ctxCli context.Context) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := GetToken(r)

			var userid = 0
			err := poolCli.QueryRow(ctxCli,
				"SELECT id FROM users WHERE token = $1", token).Scan(&userid)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if userid == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), userIdContextKey, &userid)
			r = r.WithContext(ctx)
			log.Println(ctx)
			handler.ServeHTTP(w, r)
		})
	}
}

func IsRole(requestRole string, poolCli *pgxpool.Pool, ctxCli context.Context) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := GetToken(r)

			var userRole string
			err := poolCli.QueryRow(ctxCli,
				"SELECT role FROM users WHERE token = $1", token).Scan(&userRole)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if userRole != requestRole {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
