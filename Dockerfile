FROM golang:latest

WORKDIR ./app

COPY . .

RUN go get && go install

EXPOSE 3000

CMD ["go-web-service"]