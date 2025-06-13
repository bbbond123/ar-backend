#!/bin/bash

# 设置颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 开始启动数据库服务...${NC}"

# 检查 docker-compose.db.yml 文件是否存在
if [ ! -f "docker-compose.db.yml" ]; then
    echo -e "${RED}❌ 错误: docker-compose.db.yml 文件不存在${NC}"
    exit 1
fi

# 检查 .env 文件是否存在，如果不存在则创建
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️ .env 文件不存在，创建默认配置...${NC}"
    cat > .env << EOL
BLUEPRINT_DB_USERNAME=myuser
BLUEPRINT_DB_PASSWORD=mypassword
BLUEPRINT_DB_DATABASE=mydatabase
BLUEPRINT_DB_PORT=5432
EOL
    echo -e "${GREEN}✅ 已创建默认 .env 文件${NC}"
fi

# 停止并删除现有容器
echo -e "${YELLOW}🔄 停止现有容器...${NC}"
docker-compose -f docker-compose.db.yml down

# 启动数据库服务
echo -e "${YELLOW}📦 启动数据库服务...${NC}"
docker-compose -f docker-compose.db.yml up -d

# 等待数据库启动
echo -e "${YELLOW}⏳ 等待数据库启动...${NC}"
sleep 10

# 检查数据库是否成功启动
if docker-compose -f docker-compose.db.yml ps | grep -q "Up"; then
    echo -e "${GREEN}✅ 数据库服务启动成功!${NC}"
    echo -e "${GREEN}📊 数据库信息:${NC}"
    echo -e "  - 主机: localhost"
    echo -e "  - 端口: 5432"
    echo -e "  - 用户名: ${BLUEPRINT_DB_USERNAME:-myuser}"
    echo -e "  - 数据库: ${BLUEPRINT_DB_DATABASE:-mydatabase}"
    echo -e "  - PgAdmin: http://localhost:9080"
    echo -e "    - 邮箱: admin@ar-backend.com"
    echo -e "    - 密码: admin123"
else
    echo -e "${RED}❌ 数据库服务启动失败${NC}"
    exit 1
fi

# 检查是否有备份文件
if [ -f "backup_*.sql" ]; then
    echo -e "${YELLOW}📥 发现备份文件，是否要恢复数据？(y/n)${NC}"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        # 获取最新的备份文件
        latest_backup=$(ls -t backup_*.sql | head -n1)
        echo -e "${YELLOW}🔄 正在恢复数据从 $latest_backup...${NC}"
        
        # 恢复数据
        if docker exec -i ar-backend-postgres psql -U ${BLUEPRINT_DB_USERNAME:-myuser} -d ${BLUEPRINT_DB_DATABASE:-mydatabase} < "$latest_backup"; then
            echo -e "${GREEN}✅ 数据恢复成功!${NC}"
        else
            echo -e "${RED}❌ 数据恢复失败${NC}"
            exit 1
        fi
    fi
fi

echo -e "${GREEN}✨ 数据库服务已准备就绪!${NC}" 