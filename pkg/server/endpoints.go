package server

import (
	"net/http"
)

type Endpoints interface {
	Endpoints() []Endpoint
}

type Endpoint struct {
	Method  string
	Path    string
	Handler http.Handler
}
