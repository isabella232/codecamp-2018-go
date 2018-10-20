# codecamp-2018-go

This application was made for CodeCamp NYC 2018 to demonstrate gRPC, Docker and Go. Have fun!

It's not meant to be idiomatic Go code, but rather to demonstrate a simple microservices
environment.

Slides: https://www.slideshare.net/MartinKess/building-services-with-grpc-docker-and-go

## Building
Run ./gen_protos.sh first to generate the required files. You only have to do this when you change
your protobufs.

## Running
Just run `docker-compose up` to build and run the example services

## Using the gRPC CLI 
Try using the `namely/grpc-cli` image to call your service

Create some aliases to make calling the CLI easier (substitute `docker.for.win.localhost` on Windows).
```
$ alias company_call='docker run -v `pwd`/protos/company:/defs --rm -it namely/grpc-cli call docker.for.mac.localhost:50051'

$ alias employee_call='docker run -v `pwd`/protos/employee:/defs --rm -it namely/grpc-cli call docker.for.mac.localhost:50052'
```

Create a Company
```
$ company_call CompanyService.CreateCompany "" --protofiles=company.proto

company_uuid: "3ac4f180-9410-467f-92b7-06763db0a8f1"
```

Create an Employee at that Company
```
$ employee_call EmployeeService.CreateEmployee \
  "employee: {name: 'Martin', company_uuid: '3ac4f180-9410-467f-92b7-06763db0a8f1'}" \
  --protofiles=employee.proto

employee_uuid: "10b286b2-247a-4864-afe5-f56163681af6"
company_uuid: "3ac4f180-9410-467f-92b7-06763db0a8f1"
name: "Martin"
```
List the Employees at that Company (Try running the above command a bunch of times first)
```
$ employee_call EmployeeService.ListEmployees \
    "company_uuid: '3ac4f180-9410-467f-92b7-06763db0a8f1'" --protofiles=employee.proto

employees {
  employee_uuid: "3701a099-7c67-40b2-bfac-c0e627efd0f7"
  company_uuid: "3ac4f180-9410-467f-92b7-06763db0a8f1"
  name: "Martin"
}
```
