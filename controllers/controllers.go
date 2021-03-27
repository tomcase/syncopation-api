package controllers

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type Registrar interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

func RegisterHandlers(r Registrar, prefix string) {
	registerHealthController(r, path.Join(prefix, "health"))
}
