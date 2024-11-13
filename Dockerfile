FROM ubuntu:20.04

# 避免交互式安装
ENV DEBIAN_FRONTEND=noninteractive

# 替换为阿里云的 apt 源
RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|http://mirrors.aliyun.com/ubuntu/|g' /etc/apt/sources.list

# 设置 Go Proxy 为阿里云的源
ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

# 安装 GO
RUN apt update
RUN apt install -y wget ca-certificates
RUN apt clean

# 设置 golang 的版本
ARG GO_VERSION=1.23.3

# 宿主机架构
ARG GO_ARCH

# 根据宿主机的架构选择下载不同的包
# Dockerfile 中的 RUN 命令是独立的，每个 RUN 都会在单独的 shell 中执行。
# 因此，拆分成两个 RUN 后，第一个 RUN 中设置的 GO_ARCH 变量在第二个 RUN 中不会保留，
# 这会导致第二个 RUN 无法正确解析 ${GO_ARCH} 变量。
RUN ARCH=$(uname -m) && \
    case "$ARCH" in \
        "x86_64") GO_ARCH="amd64";; \
        "aarch64") GO_ARCH="arm64";; \
        "armv7l") GO_ARCH="armv6l";; \
        *) echo "Unsupported architecture: $ARCH" && exit 1;; \
    esac && \
    wget https://mirrors.aliyun.com/golang/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
    tar -C /usr/local -xvzf go${GO_VERSION}.linux-${GO_ARCH}.tar.gz && \
    rm go${GO_VERSION}.linux-${GO_ARCH}.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# 设置工作目录
WORKDIR /app

# 将当前目录内容复制到工作目录
COPY . .

# 下载 Go 依赖并清理缓存
RUN go mod tidy

# 安装 gRPC 相关工具
RUN go get google.golang.org/grpc@latest
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 编译 Go 应用程序
RUN go build -o cache-server cache-server.go

# 导入 GOPATH/bin 路径，以便系统找到 protoc 插件
ENV PATH="$PATH:$(go env GOPATH)/bin"

# 暴露应用使用的端口
EXPOSE 8080

# 使用 ENTRYPOINT 启动，并将 ID 环境变量传递给应用
ENTRYPOINT ["sh", "-c", "./cache-server --id=$ID"]
