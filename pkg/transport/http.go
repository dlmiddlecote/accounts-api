package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	account "github.com/dlmiddlecote/api.accounts"
)

type server struct {
	router *httprouter.Router
	logger *zap.SugaredLogger
	as     account.Service
}

// NewServer returns a HTTP server for accessing the account service.
// The server implements http.Handler.
func NewServer(logger *zap.SugaredLogger, as account.Service) *server {
	router := httprouter.New()
	router.HandleOPTIONS = true

	s := server{
		router: router,
		logger: logger,
		as:     as,
	}

	// initialise router
	s.routes()

	return &s
}

func (s *server) routes() {
	s.router.GET("/accounts/:id", s.handleGetAccount())
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//
// helpers
//

func (s *server) respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			// TODO!
			panic(err)
		}
	}
}

//
// handlers
//

func (s *server) handleGetAccount() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			// id isn't an integer, respond with an error
			s.respond(w, r, http.StatusBadRequest, nil)
			return
		}

		acc, err := s.as.Account(id)
		if err != nil {
			// TODO: Handle different types of error
			s.respond(w, r, http.StatusNotFound, nil)
			return
		}

		s.respond(w, r, http.StatusOK, acc)
	}
}
