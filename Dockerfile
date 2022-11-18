FROM golang:1.19.2-alpine as build-env

WORKDIR /go/src/github.com/alikarimii/micro-with-gokit
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY services ./services

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goapp ./cmd

FROM alpine:3
RUN mkdir /app
# Create user and set ownership and permissions as required
RUN adduser -D myuser && chown -R myuser /app
WORKDIR /app
USER myuser
COPY --from=build-env /go/src/github.com/alikarimii/micro-with-gokit/goapp .
COPY migrations ./migrations

EXPOSE 15343
ENTRYPOINT ["./goapp"]
