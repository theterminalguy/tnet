# Talent Network

Talent Network is an API only service for finding and connecting with talents and companies. As a talent, you can find companies that match your skills and interests, and as a company, you can find talents that are looking for work.

# Installation 

We rely heavily on docker to run our services. Visit the docker homepage for installation instructions.

Once you have docker and docker-compose installed, you can run the `$ make up` to bring up the services locally.

# Running Tests

All tests can be run with `$ make test`.

# Known Issues

If you run into `ent` related issues, try running `$ make ent-generate` to generate the schema and `$ make up` to bring up the services.

