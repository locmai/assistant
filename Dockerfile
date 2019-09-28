FROM golang:latest
WORKDIR /go/src/github.com/locmai/app/
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=0 /go/src/github.com/locmai/app/assistant .
COPY --from=0 /go/src/github.com/locmai/app/.env .

EXPOSE 8080

ENTRYPOINT ["/assistant"]
