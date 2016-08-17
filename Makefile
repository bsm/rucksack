PKG=$(shell glide nv)

default: vet test

test:
	go test $(PKG)

vet:
	go vet $(PKG)
