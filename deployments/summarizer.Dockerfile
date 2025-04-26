FROM golang:1.24-alpine

WORKDIR /app
COPY . .
RUN go build -o summarizer cmd/summarizer/main.go

CMD ["./summarizer"]