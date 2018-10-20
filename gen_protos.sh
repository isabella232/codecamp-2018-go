#!/bin/bash -ex

# Generates all required protofiles
docker run -v `pwd`:/defs namely/protoc-all	-f protos/company/company.proto -o employee/gen/ -l go
docker run -v `pwd`:/defs namely/protoc-all	-f protos/employee/employee.proto -o employee/gen/ -l go

docker run -v `pwd`:/defs namely/protoc-all	-f protos/company/company.proto -o company/gen/ -l go

# Generate the gateways.
docker run -v `pwd`:/defs namely/gen-grpc-gateway \
    -f protos/employee/employee.proto -o gen/employee-gw -s EmployeeService

docker run -v `pwd`:/defs namely/gen-grpc-gateway \
    -f protos/company/company.proto -o gen/company-gw -s CompanyService
