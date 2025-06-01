#!/bin/bash

# 数据库初始化脚本
echo "🚀 开始初始化数据库..."

# 检查PostgreSQL容器是否运行
if ! docker ps | grep -q ar-backend-postgres; then
    echo "❌ PostgreSQL容器未运行，请先启动数据库"
    echo "   运行: docker-compose up postgres -d"
    exit 1
fi

# 执行初始化SQL
echo "📊 执行数据库初始化脚本..."
if docker exec -i ar-backend-postgres psql -U myuser -d mydatabase -f - < scripts/init_data.sql; then
    echo "✅ 数据库初始化成功!"
    echo ""
    echo "管理员账户信息:"
    echo "  邮箱: admin@ar-backend.com"
    echo "  密码: admin123"
    echo ""
    echo "🎯 你现在可以使用这个账户登录API了"
else
    echo "❌ 数据库初始化失败"
    exit 1
fi 