PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables
GOBASE=$(shell pwd)
GOHOSTOS=$(shell go env GOOS)
GOHOSTARCH=$(shell go env GOARCH)

# Directories
THIS_DIR=.
USER_DIR=$(THIS_DIR)/user-service
BIN_DIR=$(THIS_DIR)/bin

# Host & post for srv
SRV_HOST = 127.0.0.1
SRV_PORT = 8181

## user: updates user-service binary in ./bin folder
user:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/srv $(USER_DIR)/service.go 

## user-start: starts user-service at 127.0.0.1:8181
user-start:
	$(BIN_DIR)/srv start --host $(SRV_HOST) --port $(SRV_PORT)

## consul-agent: starts consul agent in development mode winh configs
consul-start: 
	consul agent -dev -config-dir=./consul.d -node=machine

## consul-connect-user: specifies service instance and proxy registration 127.0.0.1:8181
consul-connect-user:
	consul connect proxy -sidecar-for user

## consul-connect-web: specifies service instance and proxy registration 127.0.0.1:9191
consul-connect-web:
	consul connect proxy -sidecar-for web

## intention-create: creates the intention to deny access from web to user
intention-create:
	consul intention create -deny web user

## intention-delete: deletes the intention
intention-delete:
	consul intention delete web user

## user-catalog: presents the HTTP API lists all nodes hosting a user-service
user-catalog:
	curl http://localhost:8500/v1/catalog/service/user

## user-health: presents status of user-service
user-health:
	curl 'http://localhost:8500/v1/health/service/user?passing'

## web-catalog: presents the HTTP API lists all nodes hosting a web service
web-catalog:
	curl http://localhost:8500/v1/catalog/service/web

## web-health: presents status of web-service
web-health:
	curl 'http://localhost:8500/v1/health/service/web?passing'

## consul-reload: reloads consul configuration
consul-reload:
	consul reload

## consul-stop: stops consul gracefully
consul-stop:
	consul leave

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo