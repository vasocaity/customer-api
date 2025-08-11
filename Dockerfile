FROM golang:1.24.2-alpine

WORKDIR /app
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download
COPY .env ./

CMD ["air", "-c", ".air.toml"]
