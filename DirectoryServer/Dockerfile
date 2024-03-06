FROM golang:1.22.0-alpine AS build

WORKDIR /server
COPY . /server

# Build server
RUN go mod tidy
RUN go build -o /bin/server ./main.go

# Build microservices

FROM alpine:latest

COPY --from=build /bin/* /bin/

CMD ["/bin/server", "-port", "50051"]

