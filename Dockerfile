FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git

ARG GITHUB_TOKEN

ENV GOPRIVATE="github.com/pokt-foundation/*"

ENV GITHUB_TOKEN=$GITHUB_TOKEN

RUN git config --global url."https://github.com/".insteadOf "https://github.com/"

WORKDIR /go/src/github.com/pokt-foundation

COPY . /go/src/github.com/pokt-foundation/transaction-http-db

WORKDIR /go/src/github.com/pokt-foundation/transaction-http-db
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin ./main.go

FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder /go/src/github.com/pokt-foundation/transaction-http-db/bin ./

ENTRYPOINT ["/app/bin"]
