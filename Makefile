APP_NAME=rest-template
CMD_DIR=cmd/$(APP_NAME)

.PHONY: run build test tidy

run:
	cd $(CMD_DIR) && go run main.go

build:
	cd $(CMD_DIR) && go build -o ../../bin/$(APP_NAME)

test:
	go test ./...

tidy:
	go mod tidy
