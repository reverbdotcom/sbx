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

version/SBX_CHECKSUM: version/SBX_VERSION
	@export VERSION=`cat $<` && \
		curl -sL \
			"https://github.com/reverbdotcom/sbx/releases/download/$${VERSION}/Source code.tar.gz" \
			| shasum -a 256 \
			| awk '{ print $$1 }' \
			> $@

.PHONY: version/SBX_VERSION
version/SBX_VERSION:
	@test -n "$(VERSION)" || (echo "VERSION is not set" && exit 1)
	@echo $(VERSION) > $@
