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

SBX_CHECKSUM: SBX_VERSION
	@export VERSION=`cat $<` && \
		curl -sL \
			https://github.com/reverbdotcom/sbx/releases/download/$${VERSION}/sbx-darwin-arm64.tar.gz \
			| shasum -a 256 \
			| awk '{ print $$1 }' \
			> $@

.PHONY: SBX_VERSION
SBX_VERSION:
	@test -n "$(VERSION)" || (echo "VERSION is not set" && exit 1)
	@echo $(VERSION) > $@
