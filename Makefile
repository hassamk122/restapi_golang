APP_NAME=restapi_golang
BIN_NAME=gobin
BUILD_DIR=./bin
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

run:
	@echo "Running the server"
	@go run cmd/restapi/main.go || true

deps:
	@echo "Installing dependencies"
	@go mod tidy

build:
	@echo "Building the project"
	@mkdir -p ${BUILD_DIR}
	@go build -o $(BUILD_DIR)/$(BIN_NAME) cmd/restapi/main.go
	@echo "Build complete : ${BUILD_DIR}"

clean:
	@echo "Cleaning up"
	@rm -rf ${BUILD_DIR}
	@echo "Cleanup completed"

help:
	@echo "Available commands"
	@echo " -> make run		- Run the server"
	@echo " -> make deps	- Install dependencies"
	@echo " -> make build	- Build the binary"
	@echo " -> make clean	- Cleanup the binary"
