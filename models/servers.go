package models

import "context"

type ServerService interface {
	List(context.Context) ([]*Server, error)
	Insert(c context.Context, r *Server) (*Server, error)
	Delete(c context.Context, r *Server) error
}

type Server struct {
	Id              string
	Name            string
	Host            string
	Port            int32
	User            string
	Password        string
	SourcePath      string
	DestinationPath string
}
