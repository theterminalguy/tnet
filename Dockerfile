# syntax=docker/dockerfile:1

FROM golang:1.17-buster as deps

ENV PLATFORM="docker"
ARG GITHUB_USER
ENV GITHUB_USER=$GITHUB_USER
ARG GITHUB_PERSONAL_TOKEN
ENV GITHUB_PERSONAL_TOKEN=$GITHUB_PERSONAL_TOKEN

WORKDIR /app

# RUN git config \
#     --global \
#     url."https://${GITHUB_USER}:$GITHUB_PERSONAL_TOKEN@github.com".insteadOf \
#     "https://github.com"

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
FROM debian:buster-slim as production

# TODO: this is only use for testing we would change it to prod once we are ready to go live
ENV ENV production

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates \
    # start deps needed for wkhtmltopdf
    curl \
    libxrender1 \
    libjpeg62-turbo \
    fontconfig \
    libxtst6 \
    xfonts-75dpi \
    xfonts-base \
    xz-utils && \
    # stop deps needed for wkhtmltopdf
    rm -rf /var/lib/apt/lists/*

# Download the wkhtmltopdf binary
RUN curl "https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb" -L -o "wkhtmltopdf.deb"

# Install the wkhtmltopdf binary
RUN dpkg -i wkhtmltopdf.deb

COPY --from=app-build /web /web
COPY ./public/views/ ./public/views/
COPY cmd/generate/templates/talent-profile.html cmd/generate/templates/talent-profile.html

CMD ["/web"]
