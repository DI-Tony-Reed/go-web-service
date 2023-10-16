FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY * ./

RUN go mod tidy

RUN go build
RUN go build -ldflags="-X 'main.environment=production'" -o /go-web-service

EXPOSE 8080

CMD ["/go-web-service"]