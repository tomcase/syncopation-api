package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
)

func registerHealthController(r Registrar, prefix string) {
	handlerPath := path.Join(prefix, "/")
	log.Println(fmt.Sprintf("Registering %s", handlerPath))
	r.HandleFunc(path.Join(prefix, "/"), healthCheckHandler).Methods("GET", "OPTIONS")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
