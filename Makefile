
.PHONY: deps
deps:
	go mod tidy


.PHONY: build
build:
	go build -o artifacts/svc .


.PHONY: lint
lint:
	golangci-lint run --allow-parallel-runners -v -c .golangci.yml

