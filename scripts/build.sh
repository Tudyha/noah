#!/bin/bash

# 构建脚本

set -e

echo "开始构建项目..."

# 获取版本号
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S')
COMMIT_ID=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建参数
LDFLAGS="-s -w"
LDFLAGS="${LDFLAGS} -X 'main.Version=${VERSION}'"
LDFLAGS="${LDFLAGS} -X 'main.BuildTime=${BUILD_TIME}'"
LDFLAGS="${LDFLAGS} -X 'main.CommitID=${COMMIT_ID}'"

# 输出构建信息
echo "版本: ${VERSION}"
echo "构建时间: ${BUILD_TIME}"
echo "提交ID: ${COMMIT_ID}"

# 清理旧的构建文件
echo "清理旧的构建文件..."
rm -rf dist/
mkdir -p dist

# 构建不同平台的可执行文件
PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS=${PLATFORM%/*}
    GOARCH=${PLATFORM#*/}
    OUTPUT_NAME="noah-${VERSION}-${GOOS}-${GOARCH}"
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo "构建 ${GOOS}/${GOARCH}..."
    
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "${LDFLAGS}" -o "dist/${OUTPUT_NAME}" cmd/server/main.go
    
    if [ $? -eq 0 ]; then
        echo "✓ 构建成功: ${OUTPUT_NAME}"
    else
        echo "✗ 构建失败: ${OUTPUT_NAME}"
    fi
done

# 复制配置文件
echo "复制配置文件..."
cp -r config dist/
cp -r scripts dist/

# 创建压缩包
echo "创建压缩包..."
cd dist
tar -czf "noah-${VERSION}-linux-amd64.tar.gz" noah-${VERSION}-linux-amd64 config scripts
tar -czf "noah-${VERSION}-linux-arm64.tar.gz" noah-${VERSION}-linux-arm64 config scripts
tar -czf "noah-${VERSION}-darwin-amd64.tar.gz" noah-${VERSION}-darwin-amd64 config scripts
tar -czf "noah-${VERSION}-darwin-arm64.tar.gz" noah-${VERSION}-darwin-arm64 config scripts
zip -q "noah-${VERSION}-windows-amd64.zip" noah-${VERSION}-windows-amd64.exe config scripts
cd ..

echo "构建完成！"
echo "构建文件位于 dist/ 目录"