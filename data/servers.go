package data

import (
	"context"
	"log"

	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/tomcase/syncopation-api/models"
)

func (d *Db) List(c context.Context) ([]*models.Server, error) {
	dbpool, err := d.Connect(c)
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
	var servers []*models.Server
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

func (d *Db) Insert(c context.Context, r *models.Server) (*models.Server, error) {
	dbpool, err := d.Connect(c)
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

func (d *Db) Delete(c context.Context, r *models.Server) error {
	dbpool, err := d.Connect(c)
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

func mapEntityToDto(s *Server) *models.Server {
	return &models.Server{
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
