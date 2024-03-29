FROM golang:latest
WORKDIR /go/src/github.com/locmai/app/
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build .

FROM alpine:3.10.2

RUN apk add --no-cache ca-certificates

COPY --from=0 /go/src/github.com/locmai/app/assistant .
COPY --from=0 /go/src/github.com/locmai/app/example.env ./.env

EXPOSE 8080

ENTRYPOINT ["/assistant"]
