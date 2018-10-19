#!/bin/bash -ex

# Generates all required protofiles
docker run -v `pwd`:/defs namely/protoc-all	-f protos/company.proto -o employee/gen/company -l go
docker run -v `pwd`:/defs namely/protoc-all	-f protos/employee.proto -o employee/gen/employee -l go

docker run -v `pwd`:/defs namely/protoc-all	-f protos/company.proto -o company/gen/company -l go
