.PHONY: a
a:
	go run main.go "tell me a nice joke" 2> stderr.log

.PHONY: run
run: build
	./main

.PHONY: build
build:
	go build -o main