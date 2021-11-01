# syntax=docker/dockerfile:1

FROM golang:1.17.2-alpine AS deps

WORKDIR /app

COPY . ./

RUN go generate ./ent
RUN go mod download

#-----------------BUILD-----------------
FROM deps AS app-build

ENV CGO_ENABLED 0 
ENV GOOS linux 

WORKDIR /app

RUN go build -a -installsuffix cgo -o /web cmd/web/main.go

CMD ["/web"]

#-----------------TESTS-----------------
FROM deps AS tests

RUN go get -u github.com/kyoh86/richgo

CMD ["sh", "-c", "go vet ./... ; richgo test -v ./..."]
