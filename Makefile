.PHONY: tool
tool:
	go install tool

lint:
	go install tool

	golangci-lint run \
	--verbose \
	--enable-all \
	--fix
