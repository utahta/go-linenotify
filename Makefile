
install:
	@dep ensure

test:
	@go test -v -race $$(go list ./... | grep -v "vendor")
