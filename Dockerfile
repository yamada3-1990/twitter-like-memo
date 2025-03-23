FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY backend ./backend
COPY db/schema.sql ./db/schema.sql

RUN apk add --no-cache gcc musl-dev
RUN CGO_ENABLED=1 go build -o main ./backend/cmd

EXPOSE 9000

CMD ["./main"]