package endpoints

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	account "github.com/dlmiddlecote/accounts-api"
	"github.com/dlmiddlecote/accounts-api/pkg/server"
)

type accountEndpoints struct {
	logger *zap.SugaredLogger
	as     account.Service
}

func NewAccountEndpoints(logger *zap.SugaredLogger, as account.Service) *accountEndpoints {
	return &accountEndpoints{logger, as}
}

func (s *accountEndpoints) Endpoints() []server.Endpoint {
	return []server.Endpoint{
		{"GET", "/accounts/:id", s.handleGetAccount()},
	}
}

func (s *accountEndpoints) handleGetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
		if err != nil {
			// id isn't an integer, respond with an error
			server.Respond(w, r, http.StatusBadRequest, nil)
			return
		}

		acc, err := s.as.Account(id)
		if err != nil {
			// TODO: Handle different types of error
			server.Respond(w, r, http.StatusNotFound, nil)
			return
		}

		server.Respond(w, r, http.StatusOK, acc)
	}
}
