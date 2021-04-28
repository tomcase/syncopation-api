package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomcase/syncopation-api/models"
)

type DbContext interface {
	models.ServerService
}

type Controller interface {
	Register(r Registrar, prefix string, c DbContext)
}

type Registrar interface {
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

func RegisterHandlers(r Registrar, prefix string, c DbContext) {
	health := &HealthController{name: "health"}
	sync := &SyncController{name: "sync", ctx: c}
	servers := &ServersController{name: "servers", ctx: c}

	health.Register(r, prefix)
	sync.Register(r, prefix)
	servers.Register(r, prefix)
}
