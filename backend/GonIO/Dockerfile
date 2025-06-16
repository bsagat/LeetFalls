FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o gonIO cmd/main.go

CMD ["./gonIO"]
