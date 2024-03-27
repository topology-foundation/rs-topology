default: gram

.PHONY: gram
gram:
	go build -o ./bin/gram ./cmd/gram/main.go
	@echo "Finished building. Run \"./bin/gram\" to launch gram."

.PHONY: clean
clean:
	rm -rf bin
