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

