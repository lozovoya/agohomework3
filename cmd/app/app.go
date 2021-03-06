package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lozovoya/agohomework3/cmd/app/dto"
	"github.com/lozovoya/agohomework3/cmd/app/md"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var (
	ErrServer = errors.New("internal server error")
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Account     string             `json:"account"`
	Name        string             `json:"name"`
	Suggestions []Suggestion
	Operations  []Operation
}

type Suggestion struct {
	Id    int    `json:"sugid" bson:"sugid"`
	Icon  string `json:"icon" bson:"icon"`
	Title string `json:"title" bson:"title`
	Link  string `json:"link" bson:"link"`
}

type Operation struct {
	Id          int    `json:"oppid"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

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

	s.mux.With(logMd, identMd, roleCheckerMd("USER", s.poolCli, s.ctxCli),
		authMd(s.poolCli, s.ctxCli)).Get("/payments", s.Payments)
	s.mux.With(logMd, identMd, roleCheckerMd("SERVICE", s.poolCli, s.ctxCli),
		authMd(s.poolCli, s.ctxCli)).Post("/addsuggestion", s.AddSuggestion)

	return nil
}

func (s *Server) Payments(w http.ResponseWriter, r *http.Request) {

	var userid = 0
	userid = md.GetUserId(r)
	if userid == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	err := s.dbSug.Collection("users").FindOne(s.ctxSug,
		bson.D{{"userid", userid}}).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) AddSuggestion(w http.ResponseWriter, r *http.Request) {

	var userid = 0
	userid = md.GetUserId(r)
	if userid == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var newSuggestion *dto.SuggestionDTO
	err := json.NewDecoder(r.Body).Decode(&newSuggestion)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suggestion := Suggestion{
		Id:    newSuggestion.Sugid,
		Icon:  newSuggestion.Icon,
		Title: newSuggestion.Title,
		Link:  newSuggestion.Link,
	}

	match := bson.M{"userid": newSuggestion.UserId}
	change := bson.M{"$push": bson.M{"suggestions": suggestion}}

	result, err := s.dbSug.Collection("users").UpdateOne(s.ctxSug, match, change)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
