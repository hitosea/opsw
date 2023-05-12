export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=darwin:amd64 darwin:arm64 linux:386 linux:amd64 linux:arm linux:arm64 linux:mips64 linux:mips64le

all: assets
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/opsw_$${target_suffix};\
		echo "Build $${os}-$${arch} done";\
	)
	@cp ./release/opsw_linux_arm ./release/opsw_linux_aarch
	@cp ./release/opsw_linux_arm64 ./release/opsw_linux_aarch64
	@cp ./release/opsw_linux_amd64 ./release/opsw_linux_x86_64

build: asset
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o .

run: build
	./opsw run --mode debug

asset:
	env CGO_ENABLED=0 go-assets-builder shell -o assets/shell.go -p assets -v Shell
	env CGO_ENABLED=0 go-assets-builder web/dist -o assets/web.go -p assets -v Web

clean:
	@rm -f ./opsw
	@rm -rf ./release


# 提示 go-assets-builder: No such file or directory 时解決辦法
# go get github.com/jessevdk/go-assets-builder
# go install github.com/jessevdk/go-assets-builder@latest