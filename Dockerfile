FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app /root/

EXPOSE 80

CMD ["./app"]