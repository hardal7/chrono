APP_NAME=chrono
BUILD_DIR=bin
DOCKER=docker compose

.PHONY: build run clean migrate

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

run: build
	./$(BUILD_DIR)/$(APP_NAME)

migrate:
	go run internal/migration/migration.go

lint:
	golangci-lint run

test:
	go test -v ./internal/...

clean:
	rm -rf $(BUILD_DIR)

dev-down:
	$(DOCKER) --env-file .env -f deployments/compose-dev.yml down

dev-up: dev-down
	$(DOCKER) --env-file .env -f deployments/compose-dev.yml up

test-down:
	$(DOCKER) --env-file .env.test -f deployments/compose-test.yml down

test-up: test-down
	$(DOCKER) --env-file .env.test -f deployments/compose-test.yml up
	# --exit-code-from api
