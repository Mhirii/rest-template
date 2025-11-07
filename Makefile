APP_NAME=rest-template
CMD_DIR=cmd/$(APP_NAME)

.PHONY: run build test tidy

run:
	go run $(CMD_DIR)/main.go $(ARGS)

build:
	cd $(CMD_DIR) && go build -o ../../bin/$(APP_NAME)

test:
	go test ./...

tidy:
	go mod tidy
