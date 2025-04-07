.PHONY: mock
mock:
	docker run --rm -i -v ${PWD}:/src -w /src vektra/mockery --dir=internal/app/port --output=test/mocks --all

.PHONY: lint
lint:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run --fix ./...
