#!/bin/bash

# 生产环境构建脚本

echo "🚀 开始构建生产版本..."

# 设置环境变量
export ENVIRONMENT=production

# 清理旧的构建
echo "🧹 清理旧的构建文件..."
rm -rf build/bin/*

# 构建
echo "🔨 构建中..."
wails build -clean

# 检查构建结果
if [ $? -eq 0 ]; then
    echo "✅ 生产版本构建成功！"
    echo "📦 输出目录: build/bin/"

    # 显示构建产物
    echo ""
    echo "构建产物:"
    ls -lh build/bin/
else
    echo "❌ 构建失败"
    exit 1
fi
