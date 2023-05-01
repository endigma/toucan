test:
	gotestsum -f testname

lint:
	golangci-lint run

fix:
	go mod tidy
	gofumpt -w .
	gci write --skip-generated .
