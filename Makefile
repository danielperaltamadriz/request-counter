.PHONY: precommit
precommit: build test test-race lint


.PHONY: run
run:
	go run .


.PHONY: build
build:
	go build -o build/request-counter .


.PHONY: build-run
build-run: build
	./build/request-counter


.PHONY: test
test:
	go test -v ./...


.PHONY: test-race
test-race:
	go test -race ./...


.PHONY: lint
lint:
	golangci-lint run
