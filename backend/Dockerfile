FROM golang:alpine AS builder

# 维护者
MAINTAINER Lookeke
# 环境
ENV VERSION 1.0

WORKDIR /app

# 拷贝当前目录下可以执行文件
COPY . .

RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct

RUN mkdir -p /app/bin
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change

RUN go mod download && go mod verify
RUN go env -w CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && go build -o /app/bin/ ./...

# 如果需要特别需求, 需要依赖Linux系统的, 可以使用Ubuntu, Rockyinux等Linux发行版, 这里没有任何需求, 不使用任何环境
#FROM alpine:latest

# 构建一个完全空白的 Docker 镜像，这意味着你从头开始构建镜像，而不是基于任何现有的基础镜像。
# 使用 FROM scratch 可以创建非常轻量级的镜像，因为它不包含任何操作系统或软件包，只包括你显式添加的文件和依赖项
FROM scratch

WORKDIR /app

# 从 builder 镜像中复制 Go 二进制文件到当前镜像中
COPY --from=builder /app/bin .

# 定义容器启动时运行的命令, 把<appname> 替换为go mod xxx生成的名称, 也就是go build 之后的应用名称, 例如backend
ENTRYPOINT ["/app/<appname>", "-conf", "/data/conf"]


