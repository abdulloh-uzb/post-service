CURRENT_DIR=$(shell pwd)
APP=template
APP_CMD_DIR=./cmd

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}

lint: ## Run golangci-lint with printing to stdout
	golangci-lint -c .golangci.yaml run --build-tags "musl" ./...
migrate_up:
	migrate -source file://migrations -database postgres://postgres:1234@dbpost:5434/postdb up

migrate_down:
	migrate -source file://migrations -database postgres://postgres:1234@localhost:5432/post_service down

migrate_force:
	migrate -path migrations/ -database postgres://postgres:1234@localhost:5432/post?sslmode=disable force 1