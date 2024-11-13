# Build the application from source
FROM golang:1.23-bullseye AS build-stage

COPY go.mod go.sum ./
RUN go mod download

COPY ./server ./
COPY ./.env.production ./.env.production
COPY ./main.go ./main.go

RUN GOOS=linux go build -ldflags="-X 'main.environment=production'" -o /go-web-service-production

# Deploy the application binary into a lean image
FROM debian:bullseye AS build-release-stage

WORKDIR /

COPY --from=build-stage /go-web-service-production /go-web-service-production
COPY /.env.production /.env.production

EXPOSE 8080

ENTRYPOINT ["/go-web-service-production"]