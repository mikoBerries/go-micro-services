FRONT_END_BINARY=frontApp.exe
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailerApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mail
	@echo Stopping docker images if running...
	docker-compose down
	@echo Building when required and starting docker images...
	docker-compose up --build
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo Building broker binary...
	cd ./broker-service && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ${BROKER_BINARY} ./cmd/api
	@echo Done!

## build_logger: builds the logger as a linux executable
build_logger:
	@echo Building logger binary...
	cd ./logger-service && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ${LOGGER_BINARY} ./cmd/api
	@echo Done!

## build_auth: builds the Authentication binary as a linux executable
build_auth:
	@echo Building auth binary...
	cd ./authentication-service && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ${AUTH_BINARY} ./cmd/api
	@echo Done!

## build_mail: builds mail as a linux executable
build_mail:
	@echo Building mail binary...
	cd ./mail-service && set GOOS=linux && set GOARCH=amd64 && set CGO_ENABLED=0 && go build -o ${MAIL_BINARY} ./cmd/api
	@echo Done!

## build_front: builds the frone end binary
build_front:
	@echo Building front end binary...
	cd ./front-end && set CGO_ENABLED=0 && set GOOS=windows&& go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo Done!

## start: starts the front end
start: build_front
	@echo Starting front end
	cd ./front-end && start /B ${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo Stopping front end...
	@taskkill /IM "${FRONT_END_BINARY}" /F
	@echo "Stopped front end!"

##front starting front-end
front:
	cd ./front-end/cmd/web && go run *

## docker_rm_stage rm images used for Builder Labeled by (stage=Builder)
docker-rm-stage:
	docker image prune --filter label=stage=Builder

## log using mongosh
mongosh:
	mongosh mongodb://192.168.99.100:27017/logs --username admin --authenticationDatabase admin 