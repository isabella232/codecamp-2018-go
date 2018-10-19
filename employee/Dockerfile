FROM golang:alpine AS build
RUN apk add --no-cache git
WORKDIR /go/src/github.com/namely/codecamp-2018-go/employee
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine
COPY --from=build /go/bin/employee /usr/local/bin/
CMD ["employee"]
