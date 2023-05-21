export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = opsw

VERSION			:= $(shell git describe --tags --always --match="v*" 2> /dev/null || echo v0.0.0)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)

GOCGO 	:= env CGO_ENABLED=1
LDFLAGS	:= -s -w -X "$(MODULE)/vars.Version=$(VERSION)" -X "$(MODULE)/vars.CommitSHA=$(VERSION_HASH)"

run: build
	./opsw run --mode debug

release: asset
	$(GOCGO) GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./opsw-linux-amd64/opsw
	$(GOCGO) GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc-10 go build -trimpath -ldflags "$(LDFLAGS)" -o ./opsw-linux-arm64/opsw

build: asset
	$(GOCGO) go build -trimpath -ldflags "$(LDFLAGS)" -o .

asset:
	$(GOCGO) go-assets-builder shell -o assets/shell.go -p assets -v Shell
	$(GOCGO) go-assets-builder web/dist -o assets/web.go -p assets -v Web
	$(GOCGO) go-assets-builder database/*.sql -o assets/database.go -p assets -v Database

clean:
	@rm -f ./$(MODULE)


# 提示 go-assets-builder: No such file or directory 时解決辦法
# go get github.com/jessevdk/go-assets-builder
# go install github.com/jessevdk/go-assets-builder@latest