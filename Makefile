.PHONY: default
default: test

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build .

.PHONY: build-bench
build-bench:
	go build -o fib ./bench
