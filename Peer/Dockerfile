FROM golang:1.22.0-alpine AS build

WORKDIR /peer
COPY . /peer

RUN go mod tidy
# Build peer
RUN go build -o /bin/peer ./main.go

# Build microservices
WORKDIR /peer/microservices/download/
RUN go mod tidy
RUN go build -o /bin/download ./main.go

WORKDIR /peer/microservices/list/
RUN go mod tidy
RUN go build -o /bin/list ./main.go

WORKDIR /peer/microservices/upload/
RUN go mod tidy
RUN go build -o /bin/upload ./main.go

FROM alpine:latest

COPY --from=build /bin/* /bin/

CMD ["/bin/peer", "-download", "/bin/download", "-list", "/bin/list", "-upload", "/bin/upload"]
