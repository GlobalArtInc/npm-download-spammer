FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o npm-download-spammer .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/npm-download-spammer .
COPY --from=builder /app/npm-downloads-increaser.json .

ENV NPM_PACKAGE_NAME=""
ENV NPM_NUM_DOWNLOADS=1000
ENV NPM_MAX_CONCURRENT_DOWNLOAD=300
ENV NPM_DOWNLOAD_TIMEOUT=3000

ENTRYPOINT ["sh", "-c", "/app/npm-download-spammer"]
