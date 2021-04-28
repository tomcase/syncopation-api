package postgres

import (
	"context"
	"log"

	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/tomcase/syncopation-data/postgres"
)

type ServerService interface {
	List(context.Context) ([]*ServerDto, error)
	Insert(c context.Context, r *ServerDto) (*ServerDto, error)
	Delete(c context.Context, r *ServerDto) error
}

func (*Db) List(c context.Context) ([]*ServerDto, error) {
	dbpool, err := postgres.Connect(c)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := `
		SELECT id, name, host, port, user_name, password, source, destination
		FROM servers;
	`

	rows, err := dbpool.Query(c, query)
	if err != nil {
		return nil, err
	}
	var servers []*ServerDto
	for rows.Next() {
		if err := rows.Err(); err != nil {
			log.Printf("Failed to query servers: %v\n", err)
			return nil, err
		}
		var server Server
		rows.Scan(&server.Id, &server.Name, &server.Host, &server.Port, &server.User, &server.Password, &server.SourcePath, &server.DestinationPath)
		servers = append(servers, mapEntityToDto(&server))
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return servers, nil
}

func (*Db) Insert(c context.Context, r *ServerDto) (*ServerDto, error) {
	dbpool, err := postgres.Connect(c)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	query := `
		INSERT INTO servers (name, host, port, user_name, password, source, destination)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id, name, host, port, user_name, password, source, destination;
	`

	var server Server
	err = dbpool.QueryRow(c, query, r.Name, r.Host, r.Port, r.User, r.Password, r.SourcePath, r.DestinationPath).Scan(&server.Id, &server.Name, &server.Host, &server.Port, &server.User, &server.Password, &server.SourcePath, &server.DestinationPath)
	if err != nil {
		return nil, err
	}

	return mapEntityToDto(&server), nil
}

func (*Db) Delete(c context.Context, r *ServerDto) error {
	dbpool, err := postgres.Connect(c)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	query := `
		DELETE FROM servers
		WHERE id = $1;
	`
	_, err = dbpool.Exec(c, query, r.Id)
	if err != nil {
		return err
	}
	return nil
}

type Server struct {
	Id              pgtypeuuid.UUID
	Name            string
	Host            string
	Port            int32
	User            string
	Password        string
	SourcePath      string
	DestinationPath string
}

type ServerDto struct {
	Id              string
	Name            string
	Host            string
	Port            int32
	User            string
	Password        string
	SourcePath      string
	DestinationPath string
}

func mapEntityToDto(s *Server) *ServerDto {
	return &ServerDto{
		Id:              s.Id.UUID.String(),
		Name:            s.Name,
		Host:            s.Host,
		Port:            s.Port,
		User:            s.User,
		Password:        s.Password,
		SourcePath:      s.SourcePath,
		DestinationPath: s.DestinationPath,
	}
}

func mapDtoToEntity(s *ServerDto) *Server {
	id := pgtypeuuid.UUID{}
	err := id.Set(s.Id)
	if err != nil {
		log.Panicln("Server has nil ID! PANIC")
	}
	return &Server{
		Id:              id,
		Name:            s.Name,
		Host:            s.Host,
		Port:            s.Port,
		User:            s.User,
		Password:        s.Password,
		SourcePath:      s.SourcePath,
		DestinationPath: s.DestinationPath,
	}
}
