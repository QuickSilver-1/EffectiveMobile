FROM golang:1.22.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/music

RUN go build -o /app/main .

EXPOSE 8080

CMD ["/app/main"]