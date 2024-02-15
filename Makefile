LOCAL_BIN=$(CURDIR)/bin

install_deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

#checked by pre-commit
lint:
	$(LOCAL_BIN)/golangci-lint run --fast --config .golangci.pipeline.yaml