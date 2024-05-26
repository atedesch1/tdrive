.PHONY: lint
lint:
	gofmt -s -w .

.PHONY: build
build: 
	go build ./...

.PHONY: run
run:
	go run main.go

.PHONY: debug
debug:
	DEBUG=1 go run main.go

.PHONY: test
test: build
	go test -v -cover ./...
