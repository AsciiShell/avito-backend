FROM golang:1.12 AS builder
ADD . /go/src/github.com/asciishell/avito-backend
WORKDIR /go/src/github.com/asciishell/avito-backend
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/avito-backend ./cmd/messenger-api/

FROM alpine:3.9
RUN apk add --update-cache ca-certificates
COPY --from=builder /go/src/github.com/asciishell/avito-backend/bin/avito-backend /usr/local/bin/avito-backend
WORKDIR /usr/local/etc/avito-backend
ENTRYPOINT ["/usr/local/bin/avito-backend"]
