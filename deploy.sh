#!/bin/bash

# 设置变量
APP_NAME="ar-backend"
APP_PORT="3000"

echo "🚀 开始部署 $APP_NAME..."

# 1. 构建应用
echo "📦 构建应用..."
go build -ldflags="-s -w" -o $APP_NAME .

if [ $? -ne 0 ]; then
    echo "❌ 构建失败"
    exit 1
fi

echo "✅ 构建成功"

# 2. 停止旧的服务
echo "🛑 停止旧服务..."
pkill $APP_NAME

# 等待服务完全停止
sleep 2

# 3. 启动新服务
echo "🎉 启动新服务..."
nohup ./$APP_NAME > app.log 2>&1 &

# 等待服务启动
sleep 3

# 4. 检查服务状态
echo "🔍 检查服务状态..."
if curl -f http://localhost:$APP_PORT/api/health > /dev/null 2>&1; then
    echo "✅ 服务启动成功"
    echo "📝 日志文件: app.log"
    echo "🌐 服务地址: http://localhost:$APP_PORT"
else
    echo "❌ 服务启动失败，查看日志: tail -f app.log"
    exit 1
fi 