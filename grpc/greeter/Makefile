SERVER = greeter_server
CLIENT = greeter_client

SOURCE_DIR = .

all: build

.PHONY: clean pb

pb: pb/hello.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/hello.proto

build: pb
	go build -o $(SERVER) ./server/main.go
	go build -o $(CLIENT) ./client/main.go

clean:
	go clean $(SOURCE_DIR)
	rm -f $(SERVER) $(CLIENT)
	rm -f pb/*.go
