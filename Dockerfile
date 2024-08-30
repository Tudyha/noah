# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件，并缓存依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制其余源代码
COPY . .

# 构建应用程序
RUN go build -o noah cmd/noah/main.go

# 设置执行权限
RUN chmod +x ./noah

# 环境变量
ENV SERVER_PORT=8080
ENV ADMIN_PASSWORD=123456

# 暴露端口
EXPOSE ${SERVER_PORT}

# 运行命令
CMD ["./noah"]
