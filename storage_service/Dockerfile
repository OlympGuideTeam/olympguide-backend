FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o storage_service ./cmd/main.go

EXPOSE ${STORAGE_SERVICE_PORT}
CMD ["./storage_service"]
