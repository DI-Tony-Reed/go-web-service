BINARY_NAME=go-web-service
CONTAINER_NAME=backend-development

install:
	cd frontend && npm install

watch:
	cd frontend && npm run dev

update_go_dependencies:
	go get -u ./...

test:
	docker exec $(CONTAINER_NAME) sh -c 'go test ./...'

test-coverage:
	docker exec $(CONTAINER_NAME) sh -c ' \
		go test -coverprofile=coverage.out -coverpkg=./... ./... && \
		go tool cover -func=coverage.out && \
		go tool cover -html=coverage.out -o coverage.html'
	rm server/coverage.out

build:
	docker exec $(CONTAINER_NAME) sh -c '\
		GOARCH=amd64 GOOS=darwin go build -ldflags="-X main.environment=production" -o bin/${BINARY_NAME}-darwin && \
		GOARCH=amd64 GOOS=linux go build -ldflags="-X main.environment=production" -o bin/${BINARY_NAME}-linux && \
		GOARCH=amd64 GOOS=windows go build -ldflags="-X main.environment=production" -o bin/${BINARY_NAME}-windows'

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm bin/${BINARY_NAME}-darwin
	rm bin/${BINARY_NAME}-linux
	rm bin/${BINARY_NAME}-windows
