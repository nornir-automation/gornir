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


.PHONY: save-test-example-output
save-test-example-output:
	docker-compose run gornir \
		go run /go/src/github.com/nornir-automation/gornir/examples/$(EXAMPLE)/main.go > examples/$(EXAMPLE)/output.txt

.PHONY: _test-examples
_test-examples:
	make save-test-example-output EXAMPLE=1_simple
	make save-test-example-output EXAMPLE=2_simple_with_filter
	make save-test-example-output EXAMPLE=3_grouped_simple
	make save-test-example-output EXAMPLE=4_advanced_1
	make save-test-example-output EXAMPLE=5_advanced_2
	git diff --exit-code examples/

.PHONY: test-examples
test-examples: start-dev-env _test-examples stop-dev-env
