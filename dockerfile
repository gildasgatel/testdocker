FROM golang:1.20

WORKDIR /testdocker

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY templates ./templates


RUN CGO_ENABLED=0 GOOS=linux go build -o docker-test ./cmd/server

EXPOSE 8080

CMD ["/testdocker/docker-test"]

# docker run -d -p 8080:8080 --name test-docker test-docker