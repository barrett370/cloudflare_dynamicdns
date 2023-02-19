

run: 
  go run ./main.go

build:
  CGO_ENABLED=0 go build -o bin/app .

fmt:
  go fmt ./...
  goimports -w .

lint:
  golangci-lint run
  go vet ./...
