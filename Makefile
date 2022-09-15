format: .bin/goimports
	.bin/goimports -w .

.bin/goimports:
	GOBIN=$(shell pwd)/.bin go install golang.org/x/tools/cmd/goimports@latest
