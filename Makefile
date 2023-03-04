test:
	go test ./...

fetch:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

lint: fetch
	golangci-lint run ./...

build-cli:
	GOOS=linux GOARCH=amd64 go build -o im-chatgpt main.go

