# Migrate UP

migrate -source file://migrations -database "postgres://${USER}:${PASSWORD}@postgres.home.york:5432/syncopation?sslmode=disable" up

# Migrate DOWN

migrate -source file://migrations -database "postgres://${USER}:${PASSWORD}@postgres.home.york:5432/syncopation?sslmode=disable" down 1
