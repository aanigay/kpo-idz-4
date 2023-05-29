FROM golang:1.19-alpine as builder

COPY . /auth-app

WORKDIR /auth-app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -o auth-app-service ./cmd/main.go

#**# Это контейнер с ранее собранным приложением**

FROM alpine:latest as server

WORKDIR /root/

COPY --from=builder /auth-app/auth-app-service .

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

RUN chown root:root auth-app-service

EXPOSE 8003

CMD ["./auth-app-service"]