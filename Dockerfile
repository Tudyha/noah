FROM alpine:latest

# 设置国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

# 更新系统并安装基本的构建依赖
RUN apk update && \
    apk add --no-cache \
    nodejs \
    npm \
    wget \
    tar \
    build-base \
    git \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# 设置npm镜像源
RUN npm config set registry https://registry.npmmirror.com

ARG GO_VERSION=1.25.5
ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

RUN wget https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz -O go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

# 验证 Go 安装
RUN go version

# 设置go镜像源
ENV GOPROXY=https://goproxy.cn,direct

# fix sqllite问题
RUN go env -w CGO_CFLAGS='-O2 -g -D_LARGEFILE64_SOURCE'

WORKDIR /app