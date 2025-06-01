#!/bin/bash

# =======================================================================
# Docker 化部署脚本 (合并版本)
# 
# 功能说明:
# 1. 前端构建 - 使用 Node.js/npm 构建 React 应用
# 2. 后端容器构建 - 构建 Go 应用的 Docker 镜像  
# 3. 数据库服务启动 - PostgreSQL 容器化部署
# 4. nginx 反向代理 - SSL 终端 + 静态文件服务
# 5. 健康检查 - 验证所有服务正常运行
#
# 使用方法:
# ./scripts/deploy-docker.sh              # 完整部署
# ./scripts/deploy-docker.sh --frontend   # 仅更新前端
# ./scripts/deploy-docker.sh --backend    # 仅更新后端
# =======================================================================

set -e

# =======================================================================
# 参数解析和配置
# =======================================================================

FRONTEND_ONLY=false
BACKEND_ONLY=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --frontend)
            FRONTEND_ONLY=true
            shift
            ;;
        --backend)
            BACKEND_ONLY=true
            shift
            ;;
        --help|-h)
            echo "使用方法:"
            echo "  $0              # 完整部署"
            echo "  $0 --frontend   # 仅更新前端"
            echo "  $0 --backend    # 仅更新后端"
            exit 0
            ;;
        *)
            echo "未知参数: $1"
            echo "使用 --help 查看帮助"
            exit 1
            ;;
    esac
done

echo "🚀 开始 Docker 化部署..."

# 配置变量
ENV_FILE=".env.docker"
COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"

# =======================================================================
# 环境检查
# =======================================================================

echo "🔍 检查运行环境..."

# 检查必要文件
if [ ! -f "$ENV_FILE" ]; then
    echo "❌ 错误: $ENV_FILE 文件不存在"
    echo "请复制 .env.docker 并配置正确的环境变量"
    exit 1
fi

if [ ! -f "$COMPOSE_FILE" ]; then
    echo "❌ 错误: $COMPOSE_FILE 文件不存在"
    exit 1
fi

# 检查 Docker 和 Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker 未安装"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: Docker Compose 未安装"
    exit 1
fi

# 检查 Node.js 和 npm (前端构建必需)
if ! command -v node &> /dev/null; then
    echo "❌ 错误: Node.js 未安装"
    echo "前端构建需要 Node.js，请先安装 Node.js 18+"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ 错误: npm 未安装"
    echo "前端构建需要 npm"
    exit 1
fi

echo "✅ 环境检查通过"

# =======================================================================
# 目录准备
# =======================================================================

echo "📁 创建必要的目录..."
mkdir -p certbot/letsencrypt    # SSL 证书存储
mkdir -p certbot/www           # Let's Encrypt 验证
mkdir -p nginx/logs            # nginx 日志
mkdir -p backups               # 数据库备份

echo "✅ 目录创建完成"

# =======================================================================
# 前端构建
# =======================================================================

if [ "$BACKEND_ONLY" = false ]; then
    echo ""
    echo "🔨 开始前端构建..."
    
    # 检查前端目录
    if [ ! -d "client" ]; then
        echo "❌ 错误: client 目录不存在"
        exit 1
    fi
    
    # 进入前端目录
    cd client
    
    # 检查 package.json
    if [ ! -f "package.json" ]; then
        echo "❌ 错误: package.json 文件不存在"
        exit 1
    fi
    
    # 安装或更新前端依赖
    echo "📦 检查并安装前端依赖..."
    if [ ! -d "node_modules" ]; then
        echo "首次安装前端依赖..."
        npm install
    else
        echo "更新前端依赖..."
        npm install
    fi
    
    # 构建前端应用
    echo "🏗️ 编译前端代码..."
    npm run build
    
    # 验证构建结果
    if [ ! -d "dist" ]; then
        echo "❌ 前端构建失败，未找到 dist 目录"
        echo "请检查前端构建日志"
        exit 1
    fi
    
    # 检查构建文件
    if [ ! -f "dist/index.html" ]; then
        echo "❌ 前端构建不完整，缺少 index.html"
        exit 1
    fi
    
    echo "✅ 前端构建完成"
    echo "📁 构建文件位置: ./client/dist/"
    
    # 返回根目录
    cd ..
fi

# 如果只是前端更新，直接重启 nginx 并退出
if [ "$FRONTEND_ONLY" = true ]; then
    echo ""
    echo "🔄 仅更新前端，重启 nginx..."
    
    if docker ps | grep -q ifoodme-nginx; then
        docker-compose --env-file="$ENV_FILE" restart nginx
        echo "✅ nginx 重启完成"
    else
        echo "ℹ️ nginx 容器未运行，需要完整部署"
        echo "请运行: ./scripts/deploy-docker.sh"
        exit 1
    fi
    
    echo ""
    echo "🎉 前端更新完成！"
    echo "📍 静态文件位置: ./client/dist/"
    echo "🌐 访问地址: https://www.ifoodme.com"
    exit 0
fi

# =======================================================================
# 数据库备份
# =======================================================================

if [ "$FRONTEND_ONLY" = false ]; then
    echo ""
    echo "💾 检查是否需要备份数据库..."
    
    # 检查是否存在数据库卷
    if docker volume ls | grep -q postgres_data; then
        echo "发现现有数据库，进行备份..."
        mkdir -p "$BACKUP_DIR"
        
        # 创建数据库备份
        if docker ps | grep -q ifoodme-postgres; then
            echo "数据库正在运行，创建在线备份..."
            docker-compose --env-file="$ENV_FILE" exec postgres pg_dump -U ifoodme_user ifoodme_db > "$BACKUP_DIR/database.sql"
        else
            echo "数据库未运行，创建卷备份..."
            docker run --rm -v "$(pwd)/ar-backend_postgres_data:/data" -v "$BACKUP_DIR:/backup" ubuntu tar czf /backup/postgres_data.tar.gz -C /data .
        fi
        
        echo "✅ 数据库备份完成: $BACKUP_DIR"
    else
        echo "未发现现有数据库，跳过备份"
    fi
fi

# =======================================================================
# 停止现有服务
# =======================================================================

if [ "$FRONTEND_ONLY" = false ]; then
    echo ""
    echo "🛑 停止现有容器..."
    docker-compose --env-file="$ENV_FILE" down || true
    echo "✅ 容器停止完成"
fi

# =======================================================================
# 后端构建
# =======================================================================

if [ "$FRONTEND_ONLY" = false ]; then
    echo ""
    echo "🔨 构建后端 Docker 镜像..."
    
    # 检查 Dockerfile
    if [ ! -f "Dockerfile" ]; then
        echo "❌ 错误: Dockerfile 不存在"
        exit 1
    fi
    
    # 检查 Go 项目文件
    if [ ! -f "go.mod" ]; then
        echo "❌ 错误: go.mod 不存在"
        exit 1
    fi
    
    # 构建后端镜像 (无缓存以确保最新代码)
    docker-compose --env-file="$ENV_FILE" build --no-cache ar-backend
    echo "✅ 后端镜像构建完成"
fi

# 如果只是后端更新，重启后端服务并退出
if [ "$BACKEND_ONLY" = true ]; then
    echo ""
    echo "🔄 仅更新后端，重启后端服务..."
    
    docker-compose --env-file="$ENV_FILE" up -d postgres  # 确保数据库运行
    sleep 5
    docker-compose --env-file="$ENV_FILE" up -d ar-backend
    
    echo "✅ 后端更新完成"
    exit 0
fi

# =======================================================================
# 启动服务
# =======================================================================

echo ""
echo "🚀 启动 Docker 服务..."

# 1. 首先启动数据库
echo "📊 启动数据库服务..."
docker-compose --env-file="$ENV_FILE" up -d postgres

# 等待数据库完全启动
echo "⏳ 等待数据库启动 (30秒)..."
sleep 30

# 检查数据库健康状态
echo "🔍 检查数据库状态..."
if docker-compose --env-file="$ENV_FILE" exec postgres pg_isready -U ifoodme_user &> /dev/null; then
    echo "✅ 数据库已就绪"
else
    echo "⚠️ 数据库可能未完全启动，继续等待..."
    sleep 15
fi

# 2. 启动后端服务
echo "🔧 启动后端服务..."
docker-compose --env-file="$ENV_FILE" up -d ar-backend

# 3. 启动 nginx 和其他服务
echo "🌐 启动 nginx 和辅助服务..."
docker-compose --env-file="$ENV_FILE" up -d

echo "✅ 所有服务启动完成"

# =======================================================================
# 等待服务就绪
# =======================================================================

echo ""
echo "⏳ 等待服务启动 (20秒)..."
sleep 20

# =======================================================================
# 服务状态检查
# =======================================================================

echo ""
echo "🔍 检查服务状态..."
docker-compose --env-file="$ENV_FILE" ps

# =======================================================================
# 健康检查
# =======================================================================

echo ""
echo "🏥 进行健康检查..."
sleep 10

# 检查后端健康状态
echo "🔧 检查后端 API..."
if curl -f -s http://localhost/health > /dev/null 2>&1; then
    echo "✅ 后端服务健康"
    echo "   健康检查: http://localhost/health"
else
    echo "⚠️ 后端服务可能未正常启动"
    echo "🔍 后端容器日志:"
    docker-compose --env-file="$ENV_FILE" logs --tail=20 ar-backend
fi

# 检查前端
echo "🌐 检查前端页面..."
if curl -f -s http://localhost > /dev/null 2>&1; then
    echo "✅ 前端服务健康"
    echo "   前端页面: http://localhost/"
else
    echo "⚠️ 前端服务可能未正常启动"
    echo "🔍 nginx 容器日志:"
    docker-compose --env-file="$ENV_FILE" logs --tail=20 nginx
    
    # 检查前端文件是否存在
    echo "🔍 检查前端文件:"
    if docker-compose --env-file="$ENV_FILE" exec nginx ls -la /var/www/html/ 2>/dev/null; then
        echo "前端文件已挂载"
    else
        echo "前端文件挂载可能有问题"
    fi
fi

# 检查数据库连接
echo "📊 检查数据库连接..."
if docker-compose --env-file="$ENV_FILE" exec postgres pg_isready -U ifoodme_user &> /dev/null; then
    echo "✅ 数据库连接正常"
else
    echo "⚠️ 数据库连接异常"
fi

# =======================================================================
# 部署完成信息
# =======================================================================

echo ""
echo "🎉 Docker 部署完成!"
echo ""
echo "📋 服务信息:"
echo "  🌐 网站地址: https://www.ifoodme.com"
echo "  🔧 API 健康检查: https://www.ifoodme.com/health"
echo "  📊 数据库管理: http://localhost:9080 (如果启用了 pgadmin)"
echo "  📁 前端文件: ./client/dist/"
echo ""
echo "📝 常用命令:"
echo "  查看日志: docker-compose --env-file=$ENV_FILE logs -f [service_name]"
echo "  重启服务: docker-compose --env-file=$ENV_FILE restart [service_name]"
echo "  停止服务: docker-compose --env-file=$ENV_FILE down"
echo "  查看状态: docker-compose --env-file=$ENV_FILE ps"
echo ""
echo "🔄 快速更新:"
echo "  更新前端: $0 --frontend"
echo "  更新后端: $0 --backend"
echo "  完整部署: $0"
echo ""
echo "🔐 SSL 证书配置:"
echo "  首次运行 SSL 证书: ./scripts/init-letsencrypt.sh"
echo "  手动更新证书: docker-compose --env-file=$ENV_FILE run --rm certbot renew"
echo ""

# 最终状态显示
echo "📊 当前容器状态:"
docker-compose --env-file="$ENV_FILE" ps 