.PHONY: precommit
precommit: build test test-race lint


.PHONY: build
build:
	go build -o build/request-counter .

.PHONY: run
run:
	go run .


.PHONY: build-run
build-run: build
	export TTL_SEC=20 && ./build/request-counter


.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run
