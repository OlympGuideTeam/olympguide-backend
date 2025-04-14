FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE ${API_PORT}
CMD ["./main"]
