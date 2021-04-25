package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	pbServers "github.com/tomcase/syncopation-api/protos/servers"
	"google.golang.org/grpc"
)

type ServersController struct {
	name string
}

func (c *ServersController) Register(r Registrar, prefix string) {
	handlerPath := path.Join(prefix, c.name)
	log.Println(fmt.Sprintf("Registering %s", handlerPath))
	r.HandleFunc(path.Join(handlerPath, "/"), rootHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc(path.Join(handlerPath, "/{id}"), idHandler).Methods(http.MethodDelete, http.MethodOptions)
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodDelete:
		deleteServer(w, r)
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(os.Getenv("DATA_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pbServers.NewServersClient(conn)
	_, err = c.DeleteServer(r.Context(), &pbServers.DeleteServerRequest{Id: vars["id"]})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func insertServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(os.Getenv("DATA_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	input := Server{}
	json.NewDecoder(r.Body).Decode(&input)

	c := pbServers.NewServersClient(conn)

	log.Printf("Port is %s", input.Port)
	port, err := strconv.Atoi(input.Port)
	if err != nil {
		log.Printf("Error while converting port to an int32: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := c.InsertServer(r.Context(), &pbServers.InsertServerRequest{
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

	result := response.Server

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

func listServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(os.Getenv("DATA_URL"), opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pbServers.NewServersClient(conn)

	response, err := c.ListServers(r.Context(), &pbServers.ListServersRequest{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var servers []*Server
	if response != nil {
		for _, server := range response.Servers {
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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodOptions:
		return
	case http.MethodGet:
		listServers(w, r)
	case http.MethodPost:
		insertServer(w, r)
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
