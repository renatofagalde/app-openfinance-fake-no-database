.PHONY: run build test docker docker-run clean tidy fmt

APP_NAME=app-openfinance-fake-no-database
IMAGE=$(APP_NAME):latest

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

run:
	MOCKS_FILE=./mocks.json PORT=8080 go run ./cmd/server

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/server ./cmd/server

test:
	go test ./...

docker:
	docker build -t $(IMAGE) .

docker-run:
	docker run --rm -p 8080:8080 \
	  -v $(PWD)/mocks.json:/data/mocks.json \
	  $(IMAGE)

clean:
	rm -rf bin/
