FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/messaging-app .

EXPOSE 8080

CMD ["./messaging-app"]
