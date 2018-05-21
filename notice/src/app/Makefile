BINARY_NAME=app

fmt:
	go fmt . ./...

build:
	go build -o $(BINARY_NAME) -v

run:
	./$(BINARY_NAME)
clean:
	rm -rf vendor/*
	rm -f $(BINARY_NAME)
