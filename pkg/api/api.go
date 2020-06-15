package api

import (
	"net/http"
	"strconv"

	"go.uber.org/zap"

	account "github.com/dlmiddlecote/accounts-api"
	"github.com/dlmiddlecote/kit/api"
)

type accountAPI struct {
	logger *zap.SugaredLogger
	as     account.Service
}

// NewAPI returns an implementation of api.API.
// The returned API exposes the given account service as a HTTP API.
func NewAPI(logger *zap.SugaredLogger, as account.Service) *accountAPI {
	return &accountAPI{logger, as}
}

// Endpoints implements api.API. We list all API endpoints here.
func (a *accountAPI) Endpoints() []api.Endpoint {
	return []api.Endpoint{
		// Get single account by ID
		{"GET", "/accounts/:id", a.handleGetAccount(), []api.Middleware{}},
	}
}

// handleGetAccount handles requests for retrieving a single account by ID.
func (a *accountAPI) handleGetAccount() http.Handler {
	var h http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {

		// Retrieve the ':id' param from the url path, and check it is an integer
		id, err := strconv.Atoi(api.URLParam(r, "id"))
		if err != nil {
			// id isn't an integer, respond with an error
			api.Respond(w, r, http.StatusBadRequest, nil)
			return
		}

		// Retrieve the account
		acc, err := a.as.Account(id)
		if err != nil {
			// TODO: Handle different types of error
			api.Respond(w, r, http.StatusNotFound, nil)
			return
		}

		// Return the account
		api.Respond(w, r, http.StatusOK, acc)
	}
	return h
}
