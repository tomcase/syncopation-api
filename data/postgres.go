package data

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type Db struct {
}

func Migrate() error {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"syncopation",
		driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil {
		log.Println(err)
	}
	return nil
}

func Connect(c context.Context) (*pgxpool.Pool, error) {
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	dbpool, err := pgxpool.ConnectConfig(c, dbconfig)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}

	return dbpool, nil
}
