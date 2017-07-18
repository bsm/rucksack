PKG=$(shell go list ./... | grep -v vendor)

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)
