FROM golang:1.23-alpine AS builder
WORKDIR /src

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go -o ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-app ./cmd/main.go


FROM alpine:latest
WORKDIR /app

COPY --from=builder /src/messaging-app .

EXPOSE 8080
CMD ["./messaging-app"]
