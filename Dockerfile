FROM golang:1.21

WORKDIR /app

RUN apt-get update && apt-get install -y gcc git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1

RUN go build -o food-analyzer-api ./cmd/food-analyzer-api

CMD ["./food-analyzer-api"]
