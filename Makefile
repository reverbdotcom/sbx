.PHONY: build
build:
	@go build sbx.go

.PHONY: vet
vet:
	@go vet ./...

.PHONY: fmt
fmt:
	@go fmt ./...

