FROM golang:1.24-alpine

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o two-agents .

ENV MODE=async

CMD ["./two-agents"]