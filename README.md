# Talent Network

![Talent Network](https://www.lebow.drexel.edu/sites/default/files/story/1501077931-network.jpg)

Talent Network is an API only service for finding and connecting with talents and companies. As a talent, you can find companies that match your skills and interests, and as a company, you can find talents that are looking for work.

## Architecture

```
internal
├── database
│   ├── database.go
│   ├── postgres.go
│   └── sqlite3.go
├── handler
│   ├── handler.go
│   └── user_handler.go
├── repository
│   ├── repository.go
│   └── user_repository.go
├── router
│   └── router.go
└── service
    ├── service.go
    └── user_service.go
```

The above folder structure gives an overview of the core abstraction layers of the application. Let's talk about each layer. 

### Database

This is where you define all database engines used by the application. To add a new database enigine, simply add a new file to the `database` folder and define the interface specified in the `database.go` file.

### Handler

This layer is responsible for handling all HTTP requests. It usually contains a Service and a Repository. The Service is responsible for handling the business logic and the Repository is responsible for retrieving data from the database.

### Repository

The Repository is responsible for retrieving data from the database. If you wish to query the database, you should use the Repository.

### Router

The Router maps URL paths to handlers. It also utilizes a middleware layer to handle authentication and more.

### Service

The Service is responsible for handling the business logic. It usually contains a Repository.

## Installation 

We rely heavily on docker to run our services. Visit the docker homepage for installation instructions.

Once you have docker and docker-compose installed, you can run the `$ make up` to bring up the services locally.

## Running Tests

All tests can be run with `$ make test`.

## Known Issues

If you run into `ent` related issues, try running `$ make ent-generate` to generate the schema and `$ make up` to bring up the services.

## Useful Make Commands

Run any of the following command prefixing `$ make <command>`:

```
help:  Output available commands
setup:  Builds the web container
start:  Start all services
dev:  Run the web server in dev mode without using docker
pg:  Starts the postgres server
stop:  Stop all services
destroy:  Remove all containers and images. Also, destroy all volumes
scaffold:  Generate a new resource scaffold
test:  Run all tests
hot-reload:  Enables hot reload for the web service
ent-generate:  Generate ent Assests
go-format:  Run go fmt ./... on all go files

```

## Setting up your debugger in vscode

Use the below launch.json file

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "DEBUG TenTN",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}",
            "envFile": "${workspaceFolder}/.env",
            "env": {
                "DEBUG": "true"
            }
        }
    ]
}
```

## Provisioning TLS Certificates

If you run into issues with the app starting in development mode, you can use the following command to generate a self signed certificate.

```shell
    $ openssl req -x509 -newkey rsa:4096 -keyout cert-key.pem -out cert.pem -days 365 -nodes -subj "/C=US/ST=CA/L=San Francisco/O=TenTN/CN=localhost"
```
The above will generate a self signed certificate for the domain `localhost`. However, we recommend you use [mkcert](https://github.com/FiloSottile/mkcert) to generate a certificate.

```shell
$ mkcert -key-file cert-key.pem -cert-file cert.pem localhost
```

### Turnning off TLS

If you don't want to use TLS, you can use the following command to turn it off.

```shell
$ make dev SSL=0
```

**Note:** SSL is enabled by default.


Now head over to https://localhost:`<YOUR_HTTP_PORT>` and you should be able to see the app.

## Private Repo Dependency

docker build . -t tenweb --target prod --build-arg GITHUB_USER=<your_username> --build-arg GITHUB_PERSONAL_TOKEN=<your_token>

STAGE=prod docker-compose build --build-arg GITHUB_USER=<your_username> --build-arg GITHUB_PERSONAL_TOKEN=<your_token> web


## Generating JWT Tokens locally

You can generate a token for a recruiter using the command:

```shell
$ make task name=tokgen params="type=recruiter,email=youremail@example.com"
```

To generate a token for a talent use:

```shell
$ make task name=tokgen params="type=talent,email=youremail@example.com" 
```

## Approving a New Oauth2 Client

```
$  make task name=approve-client params="<client_id>"
```

## Updating the Oauth2 Scope

```
$ make task name=update-client-scope params="client-id=<client-id>,scope=<scope1,scope2,scope3>"
```

## Querying the Internal Endpoint

Internal endpoints are only available to a few users. You will need to provide the following headers to be able to query the internal endpoints:

```
X-TN-Internal-User-ID: <your-10hl-email>
X-TN-Internal-API-Key: <api-key>
```

API key can be generated using the following command:

```shell
$ openssl rand -base64 32
```
Once generated, set the value in your .env file as INTERNAL_API_KEY
