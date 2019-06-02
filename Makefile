PROJECT="github.com/nornir-automation/gornir"

.PHONY: lint
lint:
	docker run \
		--rm \
		-v $(PWD):/go/src/$(PROJECT) \
		-w /go/src/$(PROJECT) \
		golangci/golangci-lint \
			golangci-lint run

.PHONY: start-dev-env
start-dev-env:
	docker-compose up -d

.PHONY: stop-dev-env
stop-dev-env:
	docker-compose down

.PHONY: example
example:
	docker-compose run gornir \
		go run /go/src/github.com/nornir-automation/gornir/examples/$(EXAMPLE)/main.go

.PHONY: godoc
godoc:
	docker-compose run -p 6060:6060 gornir \
		godoc -http 0.0.0.0:6060 -v
