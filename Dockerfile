FROM golang:1.23-alpine AS builder
WORKDIR /gau_authorization

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN go build -o gau-authorization-service.bin .

FROM alpine:latest
WORKDIR /gau_authorization

RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    | tar xvz -C /tmp \
    && mv /tmp/migrate /usr/local/bin/migrate \
    && chmod +x /usr/local/bin/migrate

COPY --from=builder /gau_authorization/gau-authorization-service.bin .
COPY migrations ./migrations
COPY entrypoint.sh .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]