#!/bin/bash

echo "🧹 清理 Air 和 Go 缓存..."

# 清理 Air 临时目录
if [ -d "tmp" ]; then
    echo "删除 Air 临时目录..."
    rm -rf tmp/
    echo "✅ Air 临时目录已清理"
else
    echo "ℹ️  Air 临时目录不存在"
fi

# 清理 Go 缓存
echo "清理 Go 构建缓存..."
go clean -cache
echo "✅ Go 构建缓存已清理"

# 清理编译产物
if [ -f "ar-backend" ]; then
    echo "删除编译文件..."
    rm -f ar-backend
    echo "✅ 编译文件已清理"
fi

echo "🎉 所有缓存清理完成！" 