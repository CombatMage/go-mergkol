all: test build build-windows

build:
	cd cmd/mergekotlin && go build

build-windows:
	cd cmd/mergekotlin && GOOS=windows GOARCH=amd64 go build

test:
	go test ./...
