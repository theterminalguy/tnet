# Talent Network

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
