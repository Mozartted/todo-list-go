.PHONY: run build test clean tidy


# RUN THE APP
run:
	go run cmd/api/main.go


build:
	go run -o bin/api cmd/api/main.go


test:
	go test -v ./..

clean:
	rm -rf bin/


tidy:
	go mod tidy

install:
	go mod download