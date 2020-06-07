package endpoints

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	account "github.com/dlmiddlecote/accounts-api"
	"github.com/dlmiddlecote/kit/api"
)

type accountEndpoints struct {
	logger *zap.SugaredLogger
	as     account.Service
}

func NewAccountEndpoints(logger *zap.SugaredLogger, as account.Service) *accountEndpoints {
	return &accountEndpoints{logger, as}
}

func (s *accountEndpoints) Endpoints() []api.Endpoint {
	return []api.Endpoint{
		{"GET", "/accounts/:id", s.handleGetAccount()},
	}
}

func (s *accountEndpoints) handleGetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(httprouter.ParamsFromContext(r.Context()).ByName("id"))
		if err != nil {
			// id isn't an integer, respond with an error
			api.Respond(w, r, http.StatusBadRequest, nil)
			return
		}

		acc, err := s.as.Account(id)
		if err != nil {
			// TODO: Handle different types of error
			api.Respond(w, r, http.StatusNotFound, nil)
			return
		}

		api.Respond(w, r, http.StatusOK, acc)
	}
}
