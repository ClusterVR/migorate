PKG_NAME=$(shell basename $(PWD))

test:
	go clean -testcache ./...
	go test -cover $$(go list ./... | grep -vE '(/vendor/|/mock|/gen|/_design)')

lint:
	test -z "$$(golint $$(go list ./... | grep -v /vendor/ | grep -v /mock | grep -v /gen | grep -v /_design ) | tee /dev/stderr)"

vet:
	go vet $$(go list ./... | grep -vE '(/vendor/|/mock|/gen|/_design)')

test_all: lint vet fmterr test
