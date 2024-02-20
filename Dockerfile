# Prepare node/npm
FROM node:latest AS node_base

RUN echo "NODE Version:" && node --version
RUN echo "NPM Version:" && npm --version

COPY . ./

RUN make install

# Build the application from source
FROM golang:1.22-bullseye AS build-stage

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -ldflags="-X 'main.environment=production'" -o /go-web-service-production

# Deploy the application binary into a lean image
FROM debian:bullseye AS build-release-stage

WORKDIR /

COPY --from=build-stage /go-web-service-production /go-web-service-production
COPY /.env.production /.env.production

EXPOSE 8080

ENTRYPOINT ["/go-web-service-production"]