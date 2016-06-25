PKG:=$(shell glide nv)

default: test

test:
	go test $(PKG)
