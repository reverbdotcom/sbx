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

%.test:
	@go test -v -count=1 ./$*/...

%.run:
	@go run sbx.go $*

%_test.go: always
	@go test -v -count=1 $*.go $*_test.go

.PHONY: version/SBX_VERSION
version/SBX_VERSION:
	@test -n "$(VERSION)" || (echo "VERSION is not set" && exit 1)
	@echo $(VERSION) > $@

.PHONY: always
always:

