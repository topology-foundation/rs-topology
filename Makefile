default: gram

.PHONY: gram
gram:
	go build -o ./bin/gram ./cmd/gram/main.go
	@echo "Finished building. Run \"./bin/gram\" to launch gram."
	
.PHONY: bootstrap
bootstrap:
	go build -o ./bin/bootstrap ./cmd/bootstrap/main.go
	@echo "Finished building Bootstrap node. Run \"./bin/bootstrap\" to launch the Bootstrap node."

.PHONY: clean
clean:
	rm -rf bin
