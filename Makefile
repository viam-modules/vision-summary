MOD_ARCH := $(shell uname -m)
MOD_OS := $(shell uname -s)

test:
		go test ./countclassifier/

lint:
		golangci-lint run ./countclassifier/

module.tar.gz:
	go build -a -o module ./cmd/module
	tar -czf $@ module

