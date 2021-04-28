package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tomcase/syncopation-api/postgres"
)

type ServersController struct {
	name string
	ctx  postgres.ServerService
}

func (c *ServersController) Register(r Registrar, prefix string) {
	handlerPath := path.Join(prefix, c.name)
	log.Println(fmt.Sprintf("Registering %s", handlerPath))
	r.HandleFunc(path.Join(handlerPath, "/"), func(rw http.ResponseWriter, r *http.Request) { rootHandler(rw, r, c.ctx) }).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc(path.Join(handlerPath, "/{id}"), func(rw http.ResponseWriter, r *http.Request) { idHandler(rw, r, c.ctx) }).Methods(http.MethodDelete, http.MethodOptions)
}

func idHandler(w http.ResponseWriter, r *http.Request, ctx postgres.ServerService) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodDelete:
		deleteServer(w, r, ctx)
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request, ctx postgres.ServerService) {
	vars := mux.Vars(r)

	err := ctx.Delete(r.Context(), &postgres.ServerDto{Id: vars["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func insertServer(w http.ResponseWriter, r *http.Request, ctx postgres.ServerService) {
	w.Header().Set("Content-Type", "application/json")

	input := Server{}
	json.NewDecoder(r.Body).Decode(&input)
	port, err := strconv.Atoi(input.Port)
	if err != nil {
		log.Printf("Error while converting port to an int32: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := ctx.Insert(r.Context(), &postgres.ServerDto{
		Name:            input.Name,
		Host:            input.Host,
		Port:            int32(port),
		User:            input.User,
		Password:        input.Password,
		SourcePath:      input.Source,
		DestinationPath: input.Destination,
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s := &Server{
		Id:          result.Id,
		Name:        result.Name,
		Host:        result.Host,
		Port:        strconv.Itoa(int(result.Port)),
		User:        result.User,
		Password:    result.Password,
		Source:      result.SourcePath,
		Destination: result.DestinationPath,
	}

	json.NewEncoder(w).Encode(s)
}

func listServers(w http.ResponseWriter, r *http.Request, ctx postgres.ServerService) {
	w.Header().Set("Content-Type", "application/json")

	response, err := ctx.List(r.Context())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var servers []*Server
	if response != nil {
		for _, server := range response {
			s := &Server{
				Id:          server.Id,
				Name:        server.Name,
				Host:        server.Host,
				Port:        strconv.Itoa(int(server.Port)),
				User:        server.User,
				Password:    server.Password,
				Source:      server.SourcePath,
				Destination: server.DestinationPath,
			}
			servers = append(servers, s)
		}
	}
	json.NewEncoder(w).Encode(servers)
}

func rootHandler(w http.ResponseWriter, r *http.Request, ctx postgres.ServerService) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		listServers(w, r, ctx)
	case http.MethodPost:
		insertServer(w, r, ctx)
	}
}

type Server struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
