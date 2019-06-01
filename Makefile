PROJECT="github.com/nornir-automation/gornir"

.PHONY: lint
lint:
	docker run \
		--rm \
		-v $(PWD):/go/src/$(PROJECT) \
		-w /go/src/$(PROJECT) \
		golangci/golangci-lint \
			golangci-lint run
