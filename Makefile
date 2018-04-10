.PHONY: all

all:
	go fmt ./...
	protoc -I . pasty.proto --go_out=plugins=grpc:.
	go build

install: all
	cp ./pasty ~/golang/bin/p
