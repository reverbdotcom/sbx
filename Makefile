.PHONY: build
build:
	@go build sbx.go

.PHONY: vet
vet:
	@go vet ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test
test:
	@go test -v -count=1 ./...

%.run:
	@go run sbx.go $*
