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

RUN wget https://mirrors.aliyun.com/golang/go1.23.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.23.3.linux-amd64.tar.gz && \
    rm go1.23.3.linux-amd64.tar.gz

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

# 启动命令，使用环境变量 ID 作为参数传递给应用程序
CMD ["sh", "-c", "./cache-server --id=$ID"]
