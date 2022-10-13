format: .bin/goimports .bin/ory node_modules
	.bin/ory dev headers license
	.bin/goimports -w .
	npm exec -- prettier --write .

.bin/goimports: Makefile
	GOBIN=$(shell pwd)/.bin go install golang.org/x/tools/cmd/goimports@latest

.bin/ory: Makefile
	curl https://raw.githubusercontent.com/ory/meta/master/install.sh | bash -s -- -b .bin ory v0.1.44
	touch .bin/ory

node_modules: package-lock.json
	npm ci
	touch node_modules
