FROM golang:1.18-alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/pokt-foundation

# TODO Replace `backend-go-repo-template` with repo name
COPY . /go/src/github.com/pokt-foundation/backend-go-repo-template

WORKDIR /go/src/github.com/pokt-foundation/backend-go-repo-template
RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin ./main.go

FROM alpine:3.16.0
WORKDIR /app
COPY --from=builder /go/src/github.com/pokt-foundation/backend-go-repo-template/bin ./

ENTRYPOINT ["/app/bin"]
