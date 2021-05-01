# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=api

all: test build
build: 
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/main.go
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f ./dist/$(BINARY_NAME)
run:
		$(GOCLEAN)
		$(GOBUILD) -o ./dist/$(BINARY_NAME) -v ./cmd/main.go
		./dist/$(BINARY_NAME) 
build-prod:
		$(GOBUILD) -ldflags "-s -w" -o ./dist/$(BINARY_NAME) -v ./cmd/main.go 

