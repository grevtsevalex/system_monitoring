FROM golang:1.22

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -v -o /test-app ./integration-test

CMD ["/test-app", "-port", "55555"]