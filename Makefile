

# sync dependencies
dep:
	go mod download


# run unit-tests
test:
	go clean -testcache
	go test ./...
