#!/bin/bash

# 数据库迁移脚本

set -e

echo "开始数据库迁移..."

# 检查环境变量
if [ -z "$DB_TYPE" ]; then
    echo "错误: 请设置DB_TYPE环境变量 (mysql/sqlite)"
    exit 1
fi

# MySQL迁移
if [ "$DB_TYPE" = "mysql" ]; then
    echo "使用MySQL数据库"
    
    # 检查MySQL连接
    if ! command -v mysql &> /dev/null; then
        echo "错误: 未找到mysql命令"
        exit 1
    fi
    
    # 创建数据库
    echo "创建数据库..."
    mysql -h${DB_HOST:-localhost} -P${DB_PORT:-3306} -u${DB_USER:-root} -p${DB_PASSWORD} -e "CREATE DATABASE IF NOT EXISTS ${DB_NAME:-noah} CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    
    echo "MySQL数据库迁移完成"
    
# SQLite迁移
elif [ "$DB_TYPE" = "sqlite" ]; then
    echo "使用SQLite数据库"
    
    # 创建数据目录
    mkdir -p data
    
    # SQLite不需要额外的迁移步骤，表结构会在程序启动时自动创建
    echo "SQLite数据库迁移完成"
    
else
    echo "错误: 不支持的DB_TYPE: $DB_TYPE"
    exit 1
fi

echo "数据库迁移脚本执行完成"