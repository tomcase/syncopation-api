# Syncopation

Utility to keep folders in sync across devices.

## Environment Variables

`API_PORT` - The port you want to run the api on.  
`POSTGRES_CONNECTION_STRING` - Connection string for postgres

## Migrate UP

migrate -source file://migrations -database "postgres://${USER}:${PASSWORD}@localhost:5432/syncopation_dev?sslmode=disable" up

## Migrate DOWN

migrate -source file://migrations -database "postgres://${USER}:${PASSWORD}@localhost:5432/syncopation_dev?sslmode=disable" down 1
