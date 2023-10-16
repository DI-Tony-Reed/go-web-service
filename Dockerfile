# Specifies a parent image
FROM golang:latest

# Creates an application directory to hold the source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . ./

# Installs Go dependencies
#RUN go mod download
#RUN go mod verify
#
## Builds the application while setting production variable build variable
#RUN go build -ldflags="-X 'go-web-service/src/utils.environment=production'" -o /go-web-service
#
## Specifies the executable command that runs when the container starts
#CMD ["/go-web-service"]

RUN ["/bin/bash"]