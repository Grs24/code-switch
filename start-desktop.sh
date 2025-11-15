#!/bin/bash

echo "🖥️  启动桌面端应用"
echo "================================"

# 清理可能占用的端口
echo "🧹 清理端口..."
lsof -ti:9245 | xargs kill -9 2>/dev/null || true

# 进入桌面端目录
cd code-switch

echo ""
echo "🚀 启动 Wails 开发模式..."
echo "提示: 首次启动可能需要较长时间（安装依赖、编译）"
echo ""

# 启动开发模式
wails3 dev -config ./build/config.yml
