.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/gosar-linux

.PHONY: run-linux
run-linux:
	@./bin/gosar-linux
