lint:
	golangci-lint run ./...
test:
	go test ./...
run:
	go run cmd/songshift/main.go