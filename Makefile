ARCH := $(shell uname -m)
ifeq ($(ARCH), arm64)
	ARCH := arm64
else ifeq ($(ARCH), aarch64)
	ARCH := arm64
else ifeq ($(ARCH), x86_64)
	ARCH := amd64
endif

.PHONY: all
all: mocklogs logbench tracebench

.PHONY: logbench
logbench:
	GOMODULE=on CGO_ENABLED=0 go build -o bin/logbench cmd/logbench/main.go

.PHONY: tracebench
tracebench:
	GOMODULE=on CGO_ENABLED=0 go build -o bin/tracebench cmd/tracebench/main.go

.PHONY: mocklogs
mocklogs:
	GOMODULE=on CGO_ENABLED=0 go build -o bin/mocklogs tools/mocklogs/main.go

build-for-linux:
	GOOS=linux GOARCH=$(ARCH) GOMODULE=on CGO_ENABLED=0 go build -o bin/logbench cmd/logbench/main.go
	GOOS=linux GOARCH=$(ARCH) GOMODULE=on CGO_ENABLED=0 go build -o bin/tracebench cmd/tracebench/main.go
	GOOS=linux GOARCH=$(ARCH) GOMODULE=on CGO_ENABLED=0 go build -o bin/mocklogs tools/mocklogs/main.go

build-local-test-image: build-for-linux
	@cp bin/logbench logbench
	@cp bin/tracebench tracebench
	@cp bin/mocklogs mocklogs
	docker build -t localhost:5001/logbench:latest . -f Dockerfile.local
	@rm logbench
	@rm tracebench
	@rm mocklogs

push-local-test-image: build-local-test-image
	docker push localhost:5001/logbench:latest

clean:
	rm -rf bin
	rm -rf *.parquet
	rm -rf *.json
