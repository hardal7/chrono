APP_NAME=chrono
BUILD_DIR=bin
DOCKER=docker compose

.PHONY: build run lint test clean sqlc

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

run: build
	./$(BUILD_DIR)/$(APP_NAME)

lint:
	golangci-lint run

fmt:
	gofumpt -l -w .

test:
	go test -v ./internal/...

clean:
	sudo rm -rf certs/
	rm -f avatars/*
	rm -rf $(BUILD_DIR)

sqlc:
	rm -rf internal/db/sqlc && sqlc generate

dev-down:
	$(DOCKER) --env-file .env -f deployments/compose/dev.yml down

dev-up: dev-down
	$(DOCKER) --env-file .env -f deployments/compose/dev.yml up

test-down:
	$(DOCKER) --env-file .env.test -f deployments/compose/test.yml down

test-up: test-down
	$(DOCKER) --env-file .env.test -f deployments/compose/test.yml up --exit-code-from=api --abort-on-container-failure

prod-down:
	$(DOCKER) --env-file .env -f deployments/compose/prod.yml down

prod-up: prod-down
	$(DOCKER) --env-file .env -f deployments/compose/prod.yml up
