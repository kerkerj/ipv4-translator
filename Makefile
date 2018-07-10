
build_all: build_linux build_mac

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app_linux

build_mac:
	go build -o bin/app_mac
