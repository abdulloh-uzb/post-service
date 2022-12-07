CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

run:
	go run cmd/main.go

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...
migrate_up:
	migrate -source file://migrations -database postgres://abdulloh:abdulloh@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/postdb_abdulloh?sslmode=disable up

migrate_down:
	migrate -source file://migrations -database postgres://abdulloh:abdulloh@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/postdb_abdulloh?sslmode=disable down

migrate_force:
	migrate -path migrations/ -database postgres://abdulloh:abdulloh@database-1.c9lxq3r1itbt.us-east-1.rds.amazonaws.com:5432/postdb_abdulloh?sslmode=disable force 1