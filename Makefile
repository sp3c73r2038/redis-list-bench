GO := go
all: clean build

clean:
	mkdir -p target
	rm -rf target/*

build:
	$(GO) build -v -o target/redis-bench cmd/redis-bench/main.go
