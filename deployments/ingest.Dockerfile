FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go build -o ingest cmd/ingest/main.go

CMD ["./ingest"]