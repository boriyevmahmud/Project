CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...

swag-gen:
	echo ${REGISTRY}
	swag init -g api/router.go -o api/docs

up-version:
	migrate -source file:./migrations/ -database 'postgres://postgres:1@localhost:5432/users?sslmode=disable' up

down-version:
	migrate -source file:./migrations/ -database 'postgres://postgres:1@localhost:5432/users?sslmode=disable' down

create:
	migrate create -ext sql -dir migrations -seq create_test_table