APP_NAME=chrono
BUILD_DIR=bin

.PHONY: build run clean migrate

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

run: build
	./$(BUILD_DIR)/$(APP_NAME)

migrate:
	go run internal/migration/migration.go

clean:
	rm -rf $(BUILD_DIR)
