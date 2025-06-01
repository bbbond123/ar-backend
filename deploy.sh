#!/bin/bash

# 设置变量
APP_NAME="ar-backend"
APP_PORT="3000"

echo "🚀 开始部署 $APP_NAME..."

# 检查是否需要初始化数据库
INIT_DB=false
if [ "$1" = "--init-db" ] || [ "$1" = "-i" ]; then
    INIT_DB=true
    echo "🗄️ 将初始化数据库数据"
fi

# 设置环境变量
export SERVER_PORT=$APP_PORT
export ENVIRONMENT=production

echo "📋 环境配置:"
echo "  - 应用名称: $APP_NAME"
echo "  - 服务端口: $APP_PORT"
echo "  - 环境模式: $ENVIRONMENT"
echo "  - 初始化数据库: $INIT_DB"

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

# 3. 初始化数据库 (如果需要)
if [ "$INIT_DB" = true ]; then
    echo "🗄️ 初始化数据库..."
    if [ -f "scripts/init_db.sh" ]; then
        ./scripts/init_db.sh
        if [ $? -ne 0 ]; then
            echo "❌ 数据库初始化失败"
            exit 1
        fi
    else
        echo "⚠️ 找不到数据库初始化脚本"
    fi
fi

# 4. 启动新服务
echo "🎉 启动新服务..."
nohup ./$APP_NAME > app.log 2>&1 &

# 等待服务启动
sleep 5

# 5. 检查服务状态
echo "🔍 检查服务状态..."
if curl -f http://localhost:$APP_PORT/api/health > /dev/null 2>&1; then
    echo "✅ 服务启动成功"
    echo "📝 日志文件: app.log"
    echo "🌐 服务地址: http://localhost:$APP_PORT"
    echo "📖 API文档: http://localhost:$APP_PORT/swagger/index.html"
    
    if [ "$INIT_DB" = true ]; then
        echo ""
        echo "👤 默认管理员账户:"
        echo "   邮箱: admin@ar-backend.com"
        echo "   密码: admin123"
    fi
else
    echo "❌ 服务启动失败，查看日志: tail -f app.log"
    echo "最近的日志:"
    tail -20 app.log
    exit 1
fi

echo ""
echo "🎯 部署完成! 使用方法:"
echo "   普通部署: ./deploy.sh"
echo "   带数据初始化: ./deploy.sh --init-db" 