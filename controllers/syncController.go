package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/tomcase/syncopation-api/models"
	"github.com/tomcase/syncopation-api/sync"
)

type SyncController struct {
	name string
	ctx  models.ServerService
}

func (c *SyncController) Register(r Registrar, prefix string) {
	handlerPath := path.Join(prefix, c.name)
	log.Println(fmt.Sprintf("Registering %s", handlerPath))
	r.HandleFunc(path.Join(handlerPath, "/"), func(rw http.ResponseWriter, r *http.Request) { runSync(rw, r, c.ctx) }).Methods("POST", "OPTIONS")
}

func runSync(w http.ResponseWriter, r *http.Request, ctx models.ServerService) {
	if r.Method == http.MethodOptions {
		return
	}
	w.WriteHeader(http.StatusOK)

	go sync.Go(ctx)
}
