FROM golang:alpine AS build
RUN apk add --no-cache git
WORKDIR /go/src/github.com/namely/codecamp-2018-go/company
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine
COPY --from=build /go/bin/company /usr/local/bin/
CMD ["company"]
