.PHONY: run
run: build
	./chat

.PHONY: build
build:
	go build -o chat