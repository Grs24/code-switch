#!/bin/bash

# 模拟新用户环境的脚本

set -e

echo "🧪 模拟新用户环境测试"
echo "================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查是否有现有配置
if [ -d ~/.code-switch ]; then
    echo -e "${YELLOW}⚠️  检测到现有配置${NC}"
    echo ""
    echo "请选择操作："
    echo "1) 备份现有配置并清理（推荐）"
    echo "2) 直接清理（不备份）"
    echo "3) 取消"
    read -p "请输入选项 (1/2/3): " choice
    
    case $choice in
        1)
            BACKUP_DIR=~/code-switch-backup-$(date +%Y%m%d-%H%M%S)
            echo ""
            echo "📦 创建备份: $BACKUP_DIR"
            mkdir -p "$BACKUP_DIR"
            cp -r ~/.code-switch "$BACKUP_DIR/"
            echo -e "${GREEN}✅ 备份完成${NC}"
            ;;
        2)
            echo ""
            echo "⚠️  跳过备份"
            ;;
        3)
            echo ""
            echo "❌ 已取消"
            exit 0
            ;;
        *)
            echo ""
            echo "❌ 无效选项"
            exit 1
            ;;
    esac
fi

echo ""
echo "🗑️  清理现有配置..."
rm -rf ~/.code-switch
echo -e "${GREEN}✅ Code Switch 配置已清理${NC}"

# 询问是否清理 CLI 工具配置
echo ""
read -p "是否也清理 Claude Code 和 Codex 配置？(y/N): " clean_cli

if [[ $clean_cli =~ ^[Yy]$ ]]; then
    if [ -d ~/.claude ]; then
        rm -rf ~/.claude
        echo -e "${GREEN}✅ Claude Code 配置已清理${NC}"
    fi
    
    if [ -d ~/.codex ]; then
        rm -rf ~/.codex
        echo -e "${GREEN}✅ Codex 配置已清理${NC}"
    fi
fi

echo ""
echo "================================"
echo -e "${GREEN}✅ 新用户环境准备完成${NC}"
echo "================================"
echo ""
echo "📝 接下来的步骤："
echo ""
echo "1. 启动应用："
echo "   ./start-desktop.sh"
echo ""
echo "2. 完成登录"
echo ""
echo "3. 验证自动化流程："
echo "   - 查看浏览器控制台日志"
echo "   - 查看终端日志"
echo "   - 检查供应商卡片是否自动创建"
echo ""
echo "4. 测试 CLI 工具："
echo "   claude \"写一个 Hello World\""
echo "   codex \"创建一个 React 组件\""
echo ""
echo "5. 查看实时统计"
echo ""
echo "详细测试步骤请参考："
echo "   scripts/test-new-user-experience.md"
echo ""

if [ ! -z "$BACKUP_DIR" ]; then
    echo "💾 备份位置: $BACKUP_DIR"
    echo ""
    echo "恢复备份："
    echo "   rm -rf ~/.code-switch"
    echo "   cp -r $BACKUP_DIR/.code-switch ~/"
    echo ""
fi
