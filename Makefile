GOTEST ?= go test


install:
	@dep ensure

test:
	${GOTEST} -v -race ./...
