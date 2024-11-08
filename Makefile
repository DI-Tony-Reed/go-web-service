BINARY_NAME=go-web-service

install:
	cd frontend && npm install

watch:
	cd frontend && npm run dev

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags="-X 'main.environment=production'" -o bin/${BINARY_NAME}-darwin
	GOARCH=amd64 GOOS=linux go build -ldflags="-X 'main.environment=production'" -o bin/${BINARY_NAME}-linux
	GOARCH=amd64 GOOS=windows go build -ldflags="-X 'main.environment=production'" -o bin/${BINARY_NAME}-windows

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows
