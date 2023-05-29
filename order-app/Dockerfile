FROM golang:1.19-alpine as builder

COPY . /order-app

WORKDIR /order-app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -o order-app-service ./cmd/main.go

#**# Это контейнер с ранее собранным приложением**

FROM alpine:latest as server

WORKDIR /root/

COPY --from=builder /order-app/order-app-service .

COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

RUN chown root:root order-app-service

EXPOSE 8004

CMD ["./order-app-service"]