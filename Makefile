.PHONY: generate-sql
generate-sql:
	docker run --rm -v ${PWD}:/src -w /src sqlc/sqlc generate

.PHONY: mock
mock:
	docker run --rm -i -v ${PWD}:/src -w /src vektra/mockery:v2 --dir=internal/app/port --output=test/mocks --all

.PHONY: test
test:
	docker run --rm -i -v ${PWD}:/src -w /src golang:1.24.1-alpine go test -v ./...

.PHONY: lint
lint:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run --fix ./...
