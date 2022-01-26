# syntax=docker/dockerfile:1

FROM golang:1.17-buster as deps

ENV PLATFORM="docker"
ARG GITHUB_USER
ENV GITHUB_USER=$GITHUB_USER
ARG GITHUB_PERSONAL_TOKEN
ENV GITHUB_PERSONAL_TOKEN=$GITHUB_PERSONAL_TOKEN

WORKDIR /app

RUN git config \
    --global \
    url."https://${GITHUB_USER}:$GITHUB_PERSONAL_TOKEN@github.com".insteadOf \
    "https://github.com"

COPY go.* ./
RUN go mod download

COPY . ./

#-----------------BUILD-----------------
FROM deps AS app-build

RUN go build -v -o /web cmd/web/main.go

CMD ["/web"]

#-----------------HOT-RELOAD-----------------
FROM deps AS hot-reload

WORKDIR /app
ENV CGO_ENABLED 0 
ENV GOOS linux 
COPY . .
RUN go get github.com/githubnemo/CompileDaemon
ENTRYPOINT [ "./hot-reload.sh" ]

#-----------------TESTS-----------------
FROM deps AS test

ENV ENV test

RUN go get -u github.com/kyoh86/richgo

CMD ["sh", "-c", "go vet ./... ; richgo test -v ./..."]

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim as prod

# TODO: this is only use for testing we would change it to prod once we are ready to go live
ENV ENV staging

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=app-build /web /web

CMD ["/web"]
