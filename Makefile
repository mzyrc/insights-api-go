run:
	go run *.go

test:
	go clean -testcache && go test ./...

build:
	go build -o api .