FROM golang:alpine AS builder

# 维护者
# LABEL maintainer="Your Name <your@email.com>"

# 环境
ENV VERSION 1.0

WORKDIR /app

# 拷贝当前目录下可以执行文件
COPY . .

RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct

RUN mkdir -p /app/bin
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change

RUN go mod download && go mod tidy

# 打包为特定的环境的二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 && go build -o /app/bin/ ./...

# 如果需要特别需求, 需要依赖Linux系统的, 可以使用Ubuntu, Rockyinux等Linux发行版, 这里没有任何需求, 不使用任何环境
#FROM alpine:latest

# 构建一个完全空白的 Docker 镜像，这意味着你从头开始构建镜像，而不是基于任何现有的基础镜像。
# 使用 FROM scratch 可以创建非常轻量级的镜像，因为它不包含任何操作系统或软件包，只包括你显式添加的文件和依赖项
FROM scratch

WORKDIR /app

# 从 builder 镜像中复制 Go 二进制文件到当前镜像中
COPY --from=builder /app/bin .

# 定义容器启动时运行的命令, 把 backend 替换为go.mod文件里面module的名称, 例如 backend
ENTRYPOINT ["/app/backend", "-conf", "/data/conf"]

# 执行打包
# --progress=plain: 构建过程中显示的详细信息的格式
# --no-cache: 不使用缓存
# -t: 标签, 例如 myusername/myimage:v1
# web : 构建的镜像名
# . : 当前目录所有文件
# docker build --progress=plain --no-cache -t myusername/myimage:mytag backend .
