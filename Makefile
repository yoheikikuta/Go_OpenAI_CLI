.PHONY: a
a:
	go run chat.go "tell me a nice joke" 2> stderr.log

.PHONY: run
run: build
	./chat

.PHONY: build
build:
	go build -o chat