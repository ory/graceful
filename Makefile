format: .bin/goimports node_modules
	.bin/goimports -w .
	npm exec -- prettier --write .

.bin/goimports:
	GOBIN=$(shell pwd)/.bin go install golang.org/x/tools/cmd/goimports@latest

node_modules: package-lock.json
	npm ci
	touch node_modules
