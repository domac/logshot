.PHONY: all clean

all: format test build

test:
	go test -v ./logsend

format:
	gofmt -w ./logsend ./main.go

build:
	mkdir -p builds
	# 设置交叉编译参数:
	# GOOS为目标编译系统, mac os则为 "darwin", window系列则为 "windows"
	# 生成二进制执行文件 loghub_agent , 如在windows下则为 loghub_agent.exe
	GOOS="linux" GOARCH="amd64" go build -v -o builds/logshot ./main.go

clean:
	go clean -i