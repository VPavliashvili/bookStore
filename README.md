# Description
This is a simple CRUD api written in golang. <br>
Purpose of this repository is to practice API development with Go and its standard library including new [routing enhancements](https://go.dev/blog/routing-enhancements)

## How to run
* Docker and Docker compose installed and running on your system
* Navigate to project root and run `docker compose up`
* Make API call with your favorite tool or open swagger on localhost(port can be seen and changed in appsettings.json)

## Features
Feature wise this is a very simple API<br>
main emphasis was on general tooling which is listed below
* OpenAPI support
* Health check and api info
* Containerization using docker
* Hot reaload on file change even inside docker image using [CompileDaemon](https://github.com/githubnemo/CompileDaemon)
* Structured logging with [log/slog](https://pkg.go.dev/log/slog) inside file and console
* Custom routing grouping and middlewares using [net/http](https://pkg.go.dev/net/http)

## External Libraries/Dependencies
Major point of this project is to implement whole functionality with standard library only<br>
But there was some cases where third party dependencies was neccessary
* Swagger documentation with [swaggo](https://github.com/swaggo/swag) and it's [http-swagger](https://github.com/swaggo/http-swagger)
* [pgx](https://github.com/jackc/pgx) for working with posgres database
* Uuid generation for logging with [google/uuid](https://github.com/google/uuid)
