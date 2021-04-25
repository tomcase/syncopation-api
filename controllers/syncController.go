package controllers

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/tomcase/syncopation-api/messaging"
)

type SyncController struct {
	name string
}

func (c *SyncController) Register(r Registrar, prefix string) {
	handlerPath := path.Join(prefix, c.name)
	log.Println(fmt.Sprintf("Registering %s", handlerPath))
	r.HandleFunc(path.Join(handlerPath, "/"), runSync).Methods("POST", "OPTIONS")
}

func runSync(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	w.WriteHeader(http.StatusOK)

	messaging.Sync()
}
