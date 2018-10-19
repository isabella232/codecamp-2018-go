# codecamp-2018-go

This application was made for CodeCamp NYC 2018 to demonstrate gRPC, Docker and Go. Have fun!

It's not meant to be idiomatic Go code, but rather to demonstrate a simple microservices
environment.

(TODO LINK SLIDES)

## Building
Run ./gen_protos.sh first to generate the required files. You only have to do this when you change
your protobufs.

## Running
Just run `docker-compose up` to build and run the example services

## Testing
Run `grpc_test.sh` to run a script that uses the grpc-cli container to test the endpoints.

Run `http_test.sh` to run a script that uses curl to test the gRPC gateway images.
