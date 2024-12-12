FROM golang:1.22-alpine AS builder

# 安装必要的依赖
RUN apk add --no-cache build-base musl-dev gcc libc6-compat && \
    rm -rf /var/cache/apk/*

# 设置工作目录
WORKDIR /app

# 复制文件
COPY client ./client
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY pkg ./pkg
COPY web/dist ./web/dist
COPY go.mod go.sum ./

RUN GOPROXY=https://goproxy.cn,direct go mod download

# 构建应用程序
RUN CGO_ENABLED=1 go build -o noah cmd/ns/ns.go

FROM golang:1.22-alpine
WORKDIR /app

# 复制 Go 应用到最终镜像
COPY --from=builder /app /app

# 设置执行权限
RUN chmod +x ./noah

# 缓存client 依赖
RUN cd client && GOPROXY=https://goproxy.cn,direct go mod download

# 暴露端口
EXPOSE 8080

ENV HOST=127.0.0.1

# 运行
CMD ["sh", "-c", "./noah"]