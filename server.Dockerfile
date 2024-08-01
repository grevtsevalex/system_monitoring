FROM golang:1.22

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -v -o /test-app ./cmd/server

EXPOSE 8080

CMD ["/test-app", "-config", "./integration-test/server/config.toml", "-port", "55555"]