package service

import (
	"go.uber.org/zap"

	account "github.com/dlmiddlecote/accounts-api"
)

// service represents an implementation of the account service.
type service struct {
	logger *zap.SugaredLogger
}

// NewService returns a new service.
func NewService(logger *zap.SugaredLogger) *service {
	return &service{logger}
}

// Account returns the specified account if it can be found.
func (s *service) Account(id int) (*account.Account, error) {
	return &account.Account{
		ID:    id,
		Hash:  "hash",
		Title: "Title",
		URL:   "https://google.com",
	}, nil
}
