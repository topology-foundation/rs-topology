default: gram

.PHONY: gram clean build-proto lint-check

gram:
	go mod tidy
	go build -o ./bin/gram ./cmd/main.go
	@echo "Finished building. Run \"./bin/gram\" to launch gram."

clean:
	rm -rf bin

build-proto:
	buf generate proto

lint-check:
	golangci-lint run
