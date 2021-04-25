package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Controller interface {
	Register(r Registrar, prefix string)
}

type Registrar interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

func RegisterHandlers(r Registrar, prefix string) {
	health := &HealthController{name: "health"}
	sync := &SyncController{name: "sync"}
	servers := &ServersController{name: "servers"}

	health.Register(r, prefix)
	sync.Register(r, prefix)
	servers.Register(r, prefix)
}
