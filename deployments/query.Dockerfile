FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go build -o query cmd/query/main.go

CMD ["./query"]