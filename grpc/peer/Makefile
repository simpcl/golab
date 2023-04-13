BINARY = peer

SOURCE_DIR = .

all: build

.PHONY: pb clean

pb: pb/heartbeat.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/heartbeat.proto

build: pb
	go build 

clean:
	go clean $(SOURCE_DIR)
	rm -f $(BINARY)
	rm -f pb/*.go
