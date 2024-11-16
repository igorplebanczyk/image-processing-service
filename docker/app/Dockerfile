FROM golang:1.23.3 AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app /root/

ARG PORT
ENV PORT=${PORT}

EXPOSE ${PORT}

CMD ["./app"]