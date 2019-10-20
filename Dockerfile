FROM golang:1.13.1-alpine3.10 AS build

RUN apk add --no-cache git gcc
WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app
RUN mkdir -p /app/bin && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/ ./...

FROM scratch
COPY --from=build /app/bin/ /usr/bin/
