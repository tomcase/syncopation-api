FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN	go build -ldflags "-s -w" -o ./dist/app -v ./cmd/main.go 

CMD ["dist/app"]