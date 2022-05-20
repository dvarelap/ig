PWD = $(shell pwd)

# Run tests
.PHONY: test
test:
	go test $(PWD)/... --parallel=5 -coverprofile=cover.out
