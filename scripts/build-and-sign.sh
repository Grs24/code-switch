#!/bin/bash

# 构建并打包 macOS 应用（已包含签名）
# 使用方法: ./scripts/build-and-sign.sh

set -e

echo "🔨 开始构建和打包应用..."

# 使用 go-task 打包（会自动构建、打包并签名）
go run github.com/go-task/task/v3/cmd/task@latest package

echo "✅ 构建和签名完成"

# 查找生成的 .app 文件
APP_PATH="bin/CodeSwitch.app"

if [ ! -d "$APP_PATH" ]; then
    echo "❌ 找不到应用: $APP_PATH"
    exit 1
fi

# 验证签名
echo "🔍 验证签名..."
codesign --verify --verbose "$APP_PATH"

echo ""
echo "📦 创建分发包..."

# 创建 ZIP 包
cd bin
zip -r -q CodeSwitch.zip CodeSwitch.app
cd ..

ZIP_PATH="bin/CodeSwitch.zip"
ZIP_SIZE=$(du -h "$ZIP_PATH" | cut -f1)

echo ""
echo "✨ 完成！"
echo "   应用位置: $APP_PATH"
echo "   分发包: $ZIP_PATH ($ZIP_SIZE)"
echo ""
echo "💡 分享给测试人员："
echo "   1. 发送 CodeSwitch.zip"
echo "   2. 告诉他们解压后右键点击 CodeSwitch.app 选择「打开」"
