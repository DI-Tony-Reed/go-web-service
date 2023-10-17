BINARY_NAME=go-web-service

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-X 'main.environment=production'" -o ${BINARY_NAME}-darwin
	GOARCH=amd64 GOOS=linux go build -ldflags="-X 'main.environment=production'" -o ${BINARY_NAME}-linux
	GOARCH=amd64 GOOS=windows go build -ldflags="-X 'main.environment=production'" -o ${BINARY_NAME}-windows

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows