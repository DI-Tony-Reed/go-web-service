FROM golang:latest

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY * ./

RUN go mod tidy

RUN go build -ldflags="-X 'main.environment=production'" -o /go-web-service

EXPOSE 8082

CMD ["/go-web-service"]