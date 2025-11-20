#!/bin/bash

echo "🖥️  启动桌面端应用"
echo "================================"

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo "📂 当前目录: $(pwd)"
echo ""

# 清理可能占用的端口和进程
echo "🧹 清理端口和进程..."
lsof -ti:9245 | xargs kill -9 2>/dev/null || true
pkill -f "CodeSwitch" 2>/dev/null || true
sleep 1

echo ""
echo "🚀 启动 Wails 开发模式..."
echo ""

# 启动开发模式
wails3 dev -config ./build/config.yml
