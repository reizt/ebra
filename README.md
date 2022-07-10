# Echo Basic REST API
> This is a REST API template made with Echo, a Go framework.
## Prerequisites
This project requires **Docker** and **Docker Compose**.

Install [Docker](https://docs.docker.com/get-docker/), [Docker Compose](https://docs.docker.jp/compose/install.html)
```sh
$ docker -v && docker-compose -v
Docker version 20.10.16, build aa7e414
Docker Compose version 2.6.0
```
## Installation
```sh
$ git clone https://github.com/reizt/ebra.git
$ cd ebra
```
## Usage
```sh
$ docker-compose up
```
## Go Libraries
- Echo: Creates API server
- Gorm: Provides ORM
- Air: Enables hot reloading
## File Structure
```
config/
  database.go # Functions to connect DB & migrate DB
controllers/
  users_controller.go # CRUD Handlers of user data operation
models/
  user.go # Defines gorm's struct and hooks
tmp/ # Air uses this directory
.air.toml # Air settings
.gitignore
docker-compose.override.yml # Mac M1 user use different MySQL image
docker-compose.yml
Dockerfile # Go container's dockerfile
go.mod
go.sum
main.go # Mount routes and starts server
```
